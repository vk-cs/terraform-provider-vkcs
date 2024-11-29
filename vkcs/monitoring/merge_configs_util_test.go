package monitoring

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestMergeSuccess(t *testing.T) {
	tests := []struct {
		name             string
		userData         string
		monitoringScript string
		want             string
	}{
		{
			name:             "merge monitoring script with shell script",
			userData:         testUserDataScript,
			monitoringScript: testMonitoringScript,
			want:             strings.TrimLeftFunc(testMergeResultShell, unicode.IsSpace),
		},
		{
			name:             "merge monitoring script with simple cloud config",
			userData:         testUserDataCloudConfigSimple,
			monitoringScript: testMonitoringScript,
			want:             strings.TrimLeftFunc(testMergeResultCloudConfigSimple, unicode.IsSpace),
		},
		{
			name:             "merge monitoring script with cloud config",
			userData:         testUserDataCloudConfig,
			monitoringScript: testMonitoringScript,
			want:             strings.TrimLeftFunc(testMergeResultCloudConfig, unicode.IsSpace),
		},
		{
			name:             "merge windows monitoring script with powershell script",
			userData:         testUserDataWinScript,
			monitoringScript: testWinMonitoringScript,
			want:             strings.TrimLeftFunc(testMergeResultWinPowershell, unicode.IsSpace),
		},
		{
			name:             "merge windows monitoring script with cloud config",
			userData:         testUserDataCloudConfig,
			monitoringScript: testWinMonitoringScript,
			want:             strings.TrimLeftFunc(testMergeResultWinCloudConfig, unicode.IsSpace),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := MergeConfigs(test.monitoringScript, test.userData)
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMergeErrors(t *testing.T) {
	_, err := MergeConfigs(testMonitoringScript, testUserDataUnknownFormat)
	assert.EqualError(t, err, "only #!/bin/bash, #cloud-config, #ps1 user_data formats are supported, when cloud monitoring is used, given: #include")

	_, err = MergeConfigs(testWinMonitoringScript, testUserDataScript)
	assert.EqualError(t, err, "monitoring script has #ps1 format, but user_data has #!/bin/bash format. Windows does not have native support for the #!/bin/bash, try to rewrite your script in powershell format")

	_, err = MergeConfigs(testMonitoringScript, testUserDataWinScript)
	assert.EqualError(t, err, "monitoring script has #!/bin/bash format, but user_data has #ps1 format. Unix does not have native support for the #ps1, try to rewrite your script in bash format")
}

const testUserDataScript = `
#!/bin/bash
echo "Hello world"
`

const testUserDataCloudConfigSimple = `
#cloud-config
package_upgrade: true
packages:
  - git
final_message: "The system is up, after $UPTIME seconds"
`

const testUserDataCloudConfig = `
#cloud-config
package_upgrade: true
packages:
  - git
runcmd:
  - echo "Hello, world!" > /etc/motd
  - [ sh, -c, "echo 'Second command executed successfully!' >> /run/testing.txt" ]
write_files:
  - path: /etc/example_config.conf
    content: |
      [example-config]
      key=value
final_message: "The system is up, after $UPTIME seconds"
`

const testUserDataWinScript = `
#ps1
function New-Profile
{
  Write-Host "Running New-Profile function"
  $profileName = split-path $profile -leaf

  if (test-path $profile)
    {write-error "Profile $profileName already exists on this computer."}
  else
    {new-item -type file -path $profile -force }
}
`

const testUserDataUnknownFormat = `
#include
https://doc/examples/cloud-config-run-cmds.txt
https://cloud-config-boot-cmds.txt
`

const testMonitoringScript = `
#!/bin/bash

set -xueo pipefail

systemctl enable telegraf
`

const testWinMonitoringScript = `
#ps1
$installToPath = $env:ProgramFiles
Set-Log \"Success. Finish\"
`

const testMergeResultShell = `
#cloud-config
runcmd:
    - - bash
      - /run/scripts/user-script.sh
    - - bash
      - /run/scripts/cloud-monitoring-script.sh
write_files:
    - content: |-
        #!/bin/bash
        echo "Hello world"
      path: /run/scripts/user-script.sh
      permissions: "0777"
    - content: |-
        #!/bin/bash

        set -xueo pipefail

        systemctl enable telegraf
      path: /run/scripts/cloud-monitoring-script.sh
      permissions: "0777"
`

const testMergeResultCloudConfigSimple = `
#cloud-config
final_message: The system is up, after $UPTIME seconds
package_upgrade: true
packages:
    - git
runcmd:
    - - bash
      - /run/scripts/cloud-monitoring-script.sh
write_files:
    - content: |-
        #!/bin/bash

        set -xueo pipefail

        systemctl enable telegraf
      path: /run/scripts/cloud-monitoring-script.sh
      permissions: "0777"
`

const testMergeResultCloudConfig = `
#cloud-config
final_message: The system is up, after $UPTIME seconds
package_upgrade: true
packages:
    - git
runcmd:
    - echo "Hello, world!" > /etc/motd
    - - sh
      - -c
      - echo 'Second command executed successfully!' >> /run/testing.txt
    - - bash
      - /run/scripts/cloud-monitoring-script.sh
write_files:
    - content: |
        [example-config]
        key=value
      path: /etc/example_config.conf
    - content: |-
        #!/bin/bash

        set -xueo pipefail

        systemctl enable telegraf
      path: /run/scripts/cloud-monitoring-script.sh
      permissions: "0777"
`

const testMergeResultWinPowershell = `
#cloud-config
write_files:
    - content: |-
        #ps1
        function New-Profile
        {
          Write-Host "Running New-Profile function"
          $profileName = split-path $profile -leaf

          if (test-path $profile)
            {write-error "Profile $profileName already exists on this computer."}
          else
            {new-item -type file -path $profile -force }
        }
      path: C:\Program Files\Cloudbase Solutions\Cloudbase-Init\LocalScripts\user-script.ps1
      permissions: "0777"
    - content: |-
        #ps1
        $installToPath = $env:ProgramFiles
        Set-Log \"Success. Finish\"
      path: C:\Program Files\Cloudbase Solutions\Cloudbase-Init\LocalScripts\cloud-monitoring-script.ps1
      permissions: "0777"
`

const testMergeResultWinCloudConfig = `
#cloud-config
final_message: The system is up, after $UPTIME seconds
package_upgrade: true
packages:
    - git
runcmd:
    - echo "Hello, world!" > /etc/motd
    - - sh
      - -c
      - echo 'Second command executed successfully!' >> /run/testing.txt
write_files:
    - content: |
        [example-config]
        key=value
      path: /etc/example_config.conf
    - content: |-
        #ps1
        $installToPath = $env:ProgramFiles
        Set-Log \"Success. Finish\"
      path: C:\Program Files\Cloudbase Solutions\Cloudbase-Init\LocalScripts\cloud-monitoring-script.ps1
      permissions: "0777"
`
