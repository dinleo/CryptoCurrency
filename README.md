<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
</head>
<body>
<header>
    <nav>
        <h1>CLI</h1>
    </nav>
</header>
<div>
    <h4>1. HTML</h4>
    <ul>
        <li>Start : go run main.go -mode=html -port=4000</li>
    </ul>
    <h4>2. REST API</h4>
    <ul>
        <li>Start : go run main.go -mode=rest -port=4000</li>
        <li>API request : Possible API request is in api.http</li>
    </ul>
    <h4>This Code not complete yet</h4>
    <ul>
        <li>There are some hard code
            <ul>
                <li>Only BaseMiner can mining</li>
                <li>Miner Reward is 50</li>
                <li>Wallet is now just string</li>
            </ul>
        <li>Also have Bug</li>
            <ul>
                <li>You can't use your balance on mempool even if change exist</li>
                <li>Miner Reward is 50</li>
            </ul>
        <li>Functions not implemented yet</li>
        <ul>
            <li>Wallet</li>
            <li>P2P</li>
            <li>Pretty FrontEnd</li>
        </ul>
    </ul>
</div>
</body>
</html>
