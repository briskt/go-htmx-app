<?php
use Sil\PhpEnv\Env;
use Sil\Psr3Adapters\Psr3SamlLogger;

$metadata['http://idp:8106'] = [
    'name' => ['en' => 'IdP'],
    'host' => '__DEFAULT__',
    'privatekey' => 'saml.pem',
    'certificate' => 'saml.crt',
    'auth' => 'silauth',
];

// Copy configuration for port 80 and modify host.
$metadata['http://idp'] = $metadata['http://idp:8106'];
$metadata['http://idp']['host'] = 'idp';
