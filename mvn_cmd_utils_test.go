package pom_component_parsing

import (
	"strings"
	"testing"
)

func TestMvnCommandInfo_String(t *testing.T) {
	tests := []struct {
		name string
		info MvnCommandInfo
		want string
	}{
		{
			name: "Complete info",
			info: MvnCommandInfo{
				Path:             "/usr/bin/mvn",
				MvnVersion:       "3.6.3",
				UserSettingsPath: "/home/user/.m2/settings.xml",
				JavaHome:         "/usr/lib/jvm/java-11",
			},
			want: "MavenCommand: /usr/bin/mvn, JavaHome: /usr/lib/jvm/java-11, MavenVersion: 3.6.3, UserSettings: /home/user/.m2/settings.xml",
		},
		{
			name: "Empty info",
			info: MvnCommandInfo{},
			want: "MavenCommand: , JavaHome: , MavenVersion: , UserSettings: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.info.String(); got != tt.want {
				t.Errorf("MvnCommandInfo.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMvnCommandInfo_Command(t *testing.T) {
	tests := []struct {
		name         string
		info         MvnCommandInfo
		args         []string
		wantCommand  string
		wantJavaHome string
	}{
		{
			name: "Command with settings and JAVA_HOME",
			info: MvnCommandInfo{
				Path:             "/usr/bin/mvn",
				UserSettingsPath: "/home/user/.m2/settings.xml",
				JavaHome:         "/usr/lib/jvm/java-11",
			},
			args:         []string{"clean", "install"},
			wantCommand:  "/usr/bin/mvn --settings /home/user/.m2/settings.xml --batch-mode clean install",
			wantJavaHome: "/usr/lib/jvm/java-11",
		},
		{
			name: "Command without settings",
			info: MvnCommandInfo{
				Path:     "/usr/bin/mvn",
				JavaHome: "/usr/lib/jvm/java-11",
			},
			args:         []string{"package"},
			wantCommand:  "/usr/bin/mvn --batch-mode package",
			wantJavaHome: "/usr/lib/jvm/java-11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.info.Command(tt.args...)

			// 验证命令路径和参数
			gotCommand := cmd.Path + " " + joinArgs(cmd.Args[1:])
			if gotCommand != tt.wantCommand {
				t.Errorf("Command() got = %v, want %v", gotCommand, tt.wantCommand)
			}

			// 验证环境变量
			if tt.wantJavaHome != "" {
				found := false
				for _, env := range cmd.Env {
					if env == "JAVA_HOME="+tt.wantJavaHome {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Command() JAVA_HOME not found or incorrect")
				}
			}
		})
	}
}

func TestParseMvnVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "Valid version",
			input: `Apache Maven 3.6.3 (cecedd343002696d0abb50b32b541b8a6ba2883f)
Maven home: /usr/share/maven
Java version: 11.0.11, vendor: Ubuntu
Java home: /usr/lib/jvm/java-11-openjdk-amd64`,
			want: "3.6.3",
		},
		{
			name:  "Empty input",
			input: "",
			want:  "",
		},
		{
			name:  "Invalid format",
			input: "Some random text",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseMvnVersion(tt.input); got != tt.want {
				t.Errorf("parseMvnVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckMvnCommand(t *testing.T) {
	// 保存原始缓存并在测试结束后恢复
	originalCache := cachedMvnCommandResult
	defer func() {
		cachedMvnCommandResult = originalCache
	}()

	tests := []struct {
		name       string
		setupFunc  func()
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Cached result",
			setupFunc: func() {
				cachedMvnCommandResult = &_MvnCommandResult{
					rs: &MvnCommandInfo{
						Path:       "/usr/bin/mvn",
						MvnVersion: "3.6.3",
					},
					e: nil,
				}
			},
			wantErr: false,
		},
		{
			name: "Cached error",
			setupFunc: func() {
				cachedMvnCommandResult = &_MvnCommandResult{
					rs: nil,
					e:  ErrMvnNotFound,
				}
			},
			wantErr:    true,
			wantErrMsg: ErrMvnNotFound.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置测试环境
			cachedMvnCommandResult = nil
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			info, err := CheckMvnCommand()
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckMvnCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.wantErrMsg {
				t.Errorf("CheckMvnCommand() error = %v, wantErrMsg %v", err, tt.wantErrMsg)
			}

			if !tt.wantErr && info == nil {
				t.Error("CheckMvnCommand() returned nil info when no error expected")
			}
		})
	}
}

// Helper function to join command arguments
func joinArgs(args []string) string {
	return strings.Join(args, " ")
}
