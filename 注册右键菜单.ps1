param( $param1, $param2 )
# 检查并以管理员身份运行PS并带上参数
$currentWi = [Security.Principal.WindowsIdentity]::GetCurrent()
$currentWp = [Security.Principal.WindowsPrincipal]$currentWi
if( -not $currentWp.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator))
{
    $boundPara = ($MyInvocation.BoundParameters.Keys | foreach{'-{0} {1}' -f  $_ ,$MyInvocation.BoundParameters[$_]} ) -join ' '
    $currentFile = $MyInvocation.MyCommand.Definition
    $fullPara = $boundPara + ' ' + $args -join ' '
    Start-Process "$psHome\powershell.exe"   -ArgumentList "$currentFile $fullPara"   -verb runas
    return
}
# 获取当前文件路径
$currentpath = Split-Path -Parent $MyInvocation.MyCommand.Definition
$ico_path = $currentpath + "\kaf-cli.exe"
$exe_path = $currentpath + "\kaf-cli.exe" + ' "%1"'

# 注册右键菜单（使用 HKEY_CURRENT_USER，不需要管理员权限，Win10/Win11 通用）
New-Item -Force -Path Registry::HKEY_CURRENT_USER\Software\Classes\txtfile\shell\使用kaf-cli转换
New-ItemProperty -Force -Path Registry::HKEY_CURRENT_USER\Software\Classes\txtfile\shell\使用kaf-cli转换 -Name Icon -PropertyType String -Value $ico_path

New-Item -Force -Path Registry::HKEY_CURRENT_USER\Software\Classes\txtfile\shell\使用kaf-cli转换\command
New-ItemProperty -Force -Path Registry::HKEY_CURRENT_USER\Software\Classes\txtfile\shell\使用kaf-cli转换\command -Name "(default)" -PropertyType String -Value $exe_path

echo "注册右键菜单成功!"
pause
