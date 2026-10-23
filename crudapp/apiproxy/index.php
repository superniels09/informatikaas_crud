<?php
$baseURL = 'http://localhost:7000';

function isServerRunning($url) {
    $ch = curl_init($url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_NOBODY, true); // HEAD request, don't fetch body
    curl_setopt($ch, CURLOPT_TIMEOUT, 2);   // don't hang forever
    curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 2);
    curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    $error = curl_errno($ch);
    curl_close($ch);

    // If curl succeeded and we got some HTTP status back, server is up
    return $error === 0 && $httpCode > 0;
}

if (!isServerRunning($baseURL)) {
	include "../../env.php";
	global $crudRestartCmd;

	echo shell_exec('cd ../../ && ' . $crudRestartCmd);
	sleep(1);
}

// 1. Extract target path, defaulting to empty or root
$targetPath = $_GET['__target_path'] ?? '';
if ($targetPath === '') {
    http_response_code(400);
    echo json_encode(['error' => 'missing __target_path query parameter']);
    exit;
}

// 2. Remove __target_path so it isn't forwarded
$queryParams = $_GET;
unset($queryParams['__target_path']);

// 3. Rebuild the URL with remaining query parameters
$url = $baseURL . '/' . ltrim($targetPath, '/');
if (!empty($queryParams)) {
    $url .= '?' . http_build_query($queryParams);
}

// 4. Get the incoming HTTP method
$method = $_SERVER['REQUEST_METHOD'];

// 5. Read the raw request body (works for POST, PUT, PATCH, DELETE w/ body, etc.)
$body = file_get_contents('php://input');

// 6. Forward incoming request headers (except Host, which should match the target)
$forwardHeaders = [];
foreach (getallheaders() as $name => $value) {
    if (strcasecmp($name, 'Host') === 0) {
        continue;
    }
    if (strcasecmp($name, 'Content-Length') === 0) {
        continue; // let curl set this based on the actual body it sends
    }
    $forwardHeaders[] = "$name: $value";
}

$ch = curl_init();

$options = [
    CURLOPT_URL => $url,
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_HEADER => false,
    CURLOPT_CUSTOMREQUEST => $method,
    CURLOPT_HTTPHEADER => $forwardHeaders,
    // Intercept and relay headers from target server to client
    CURLOPT_HEADERFUNCTION => function ($ch, $header) {
        $length = strlen($header);
        $trimmed = trim($header);
        if ($trimmed !== '' && stripos($trimmed, 'Transfer-Encoding:') === false) {
            header($trimmed, false);
        }
        return $length;
    }
];

// 7. Attach body for methods that typically carry one
if (in_array($method, ['POST', 'PUT', 'PATCH', 'DELETE'], true) && $body !== '') {
    $options[CURLOPT_POSTFIELDS] = $body;
}

curl_setopt_array($ch, $options);

$response = curl_exec($ch);

if (curl_errno($ch)) {
    http_response_code(502);
    echo json_encode(['error' => curl_error($ch)]);
    curl_close($ch);
    exit;
}

$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
http_response_code($httpCode);
curl_close($ch);

echo $response;
?>
