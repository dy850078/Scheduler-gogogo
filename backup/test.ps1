# 設定 API Endpoint 與 Request Payload
$endpoint = "http://localhost:8080/schedule"
$payload = @{
    requested_cpu    = 4
    requested_memory = 81920
    requested_pool   = "default"
    dedicated        = $false
} | ConvertTo-Json

# 建立一個工作列表
$jobs = @()

# 發送 20 筆非同步請求
1..20 | ForEach-Object {
    $jobs += Start-Job -ScriptBlock {
        param($endpoint, $payload)
        Invoke-RestMethod -Uri $endpoint -Method POST -Body $payload -ContentType "application/json"
    } -ArgumentList $endpoint, $payload
}

# 等待所有工作完成並印出結果
$jobs | ForEach-Object {
    $result = Receive-Job -Job $_ -Wait -AutoRemoveJob
    Write-Output "[DEBUG] Response: $result"
}