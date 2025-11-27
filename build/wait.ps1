# Wait for a specified number of seconds and check the frontend server
param(
    [int]$Seconds = 5,
    [string]$Url = "http://localhost:9245"
)

Write-Host "Waiting for frontend to start..." -ForegroundColor Yellow

$elapsed = 0
$interval = 1

while ($elapsed -lt $Seconds) {
    try {
        $response = Invoke-WebRequest -Uri $Url -Method Head -TimeoutSec 1 -ErrorAction SilentlyContinue
        if ($response.StatusCode -eq 200) {
            Write-Host "✓ Frontend server is ready ($elapsed seconds)" -ForegroundColor Green
            exit 0
        }
    } catch {
        # Continue waiting
    }
    
    Start-Sleep -Seconds $interval
    $elapsed += $interval
    Write-Host "  Waiting... ($elapsed/$Seconds seconds)" -ForegroundColor Gray
}

Write-Host "✓ Waiting finished (up to $Seconds seconds)" -ForegroundColor Green
