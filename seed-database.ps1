# Xephyr Database Seeding Script
# This script inserts dummy data into the PostgreSQL database

$ErrorActionPreference = "Stop"

Write-Host "==============================================" -ForegroundColor Cyan
Write-Host "Xephyr Database Seeding Script" -ForegroundColor Cyan
Write-Host "==============================================" -ForegroundColor Cyan

# Database connection parameters
$env:PGHOST = "localhost"
$env:PGPORT = "5432"
$env:PGUSER = "xephyr"
$env:PGPASSWORD = "xephyr123"
$env:PGDATABASE = "xephyr"

# Check if psql is available
try {
    $psqlVersion = psql --version 2>$null
    Write-Host "Found PostgreSQL client: $psqlVersion" -ForegroundColor Green
} catch {
    Write-Host "Error: PostgreSQL client (psql) not found in PATH" -ForegroundColor Red
    Write-Host "Please install PostgreSQL or add psql to your PATH" -ForegroundColor Yellow
    exit 1
}

# Test database connection
Write-Host "`nTesting database connection..." -ForegroundColor Yellow
try {
    $result = psql -c "SELECT 1 as connection_test;" 2>&1
    if ($result -match "connection_test") {
        Write-Host "Database connection successful!" -ForegroundColor Green
    } else {
        throw "Connection failed"
    }
} catch {
    Write-Host "Error: Cannot connect to database" -ForegroundColor Red
    Write-Host "Make sure PostgreSQL is running and accessible" -ForegroundColor Yellow
    Write-Host "Connection details: $env:PGUSER@$env:PGHOST`:$env:PGPORT/$env:PGDATABASE" -ForegroundColor Gray
    exit 1
}

# Function to execute SQL file
function Invoke-SqlFile {
    param(
        [string]$FilePath,
        [string]$Description
    )
    
    Write-Host "`nExecuting: $Description..." -ForegroundColor Yellow -NoNewline
    try {
        psql -f $FilePath 2>&1 | Out-Null
        Write-Host " DONE" -ForegroundColor Green
        return $true
    } catch {
        Write-Host " FAILED" -ForegroundColor Red
        Write-Host "Error: $_" -ForegroundColor Red
        return $false
    }
}

# Execute seed scripts
$scripts = @(
    @{ File = "init-scripts/01_seed_data.sql"; Desc = "Organizations, Users, and Skills" },
    @{ File = "init-scripts/02_seed_projects_tasks.sql"; Desc = "Projects and Tasks" },
    @{ File = "init-scripts/03_seed_nudges_workload.sql"; Desc = "Nudges, Workload, and Scenarios" }
)

$successCount = 0
foreach ($script in $scripts) {
    if (Invoke-SqlFile -FilePath $script.File -Description $script.Desc) {
        $successCount++
    }
}

Write-Host "`n==============================================" -ForegroundColor Cyan
if ($successCount -eq $scripts.Count) {
    Write-Host "Database seeding completed successfully!" -ForegroundColor Green
    Write-Host "==============================================" -ForegroundColor Cyan
    
    # Show summary
    Write-Host "`nSeeded Data Summary:" -ForegroundColor Cyan
    Write-Host "-------------------" -ForegroundColor Gray
    
    $summaryQueries = @(
        @{ Query = "SELECT COUNT(*) as count FROM organizations;"; Label = "Organizations" },
        @{ Query = "SELECT COUNT(*) as count FROM users;"; Label = "Users" },
        @{ Query = "SELECT COUNT(*) as count FROM skills;"; Label = "Skills" },
        @{ Query = "SELECT COUNT(*) as count FROM projects;"; Label = "Projects" },
        @{ Query = "SELECT COUNT(*) as count FROM tasks;"; Label = "Tasks" },
        @{ Query = "SELECT COUNT(*) as count FROM nudges;"; Label = "Nudges" },
        @{ Query = "SELECT COUNT(*) as count FROM workload_entries;"; Label = "Workload Entries" },
        @{ Query = "SELECT COUNT(*) as count FROM scenarios;"; Label = "Scenarios" }
    )
    
    foreach ($item in $summaryQueries) {
        $result = psql -c $item.Query -t 2>$null
        $count = $result.Trim()
        Write-Host "$($item.Label.PadRight(20)): $count" -ForegroundColor White
    }
    
    Write-Host "`nYour dashboard should now display real data!" -ForegroundColor Green
    Write-Host "Start the backend: cd backend && go run cmd/server/main.go" -ForegroundColor Yellow
    Write-Host "Start the frontend: cd frontend && npm run dev" -ForegroundColor Yellow
} else {
    Write-Host "Database seeding completed with errors!" -ForegroundColor Red
    Write-Host "$successCount/$($scripts.Count) scripts executed successfully" -ForegroundColor Yellow
}
Write-Host "==============================================" -ForegroundColor Cyan
