package ssh2

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/gliderlabs/ssh"
)

func (ssh2 *SSH2) setupHandler(server *ssh.Server) {
	server.Handle(ssh2.handleConnection)
}

const welcomeMessage = `
__        ___         _  _____   ____  _     _   _ _ _
\ \      / (_)___ ___| |/ /_ _| |  _ \(_)___| |_(_) | | ___ _ __ _   _
 \ \ /\ / /| / __/ __| ' / | |  | | | | / __| __| | | |/ _ \ '__| | | |
  \ V  V / | \__ \__ \ . \ | |  | |_| | \__ \ |_| | | |  __/ |  | |_| |
   \_/\_/  |_|___/___/_|\_\___| |____/|_|___/\__|_|_|_|\___|_|   \__, |
                                                                 |___/

Welcome to the WissKI SSH Server.
You've successfully authenticated, but we don't provide shell access to
the main server. You may use this connection as part of a proxy jump to
connect to your WissKI Instance.

To connect to a WissKI named ${SLUG} you may use:

ssh -J ${DOMAIN}:${PORT} www-data@${HOSTNAME}

For more details see:

${HELP_URL}

Press CTRL-C to close this connection.
`

func (ssh2 *SSH2) handleConnection(session ssh.Session) {
	slug, _ := getAnyPermission(session.Context())

	banner := welcomeMessage
	for _, oldnew := range [][2]string{
		{"${SLUG}", slug},
		{"${HOSTNAME}", slug + "." + ssh2.Config.HTTP.PrimaryDomain},

		{"${DOMAIN}", ssh2.Config.HTTP.PanelDomain()},
		{"${PORT}", strconv.FormatUint(uint64(ssh2.Config.Listen.SSHPort), 10)},

		{"${HELP_URL}", ssh2.Config.HTTP.JoinPath("user", "ssh").String()},
	} {
		banner = strings.ReplaceAll(banner, oldnew[0], oldnew[1])
	}

	io.WriteString(session, banner)

	// wait until the user closes
	buffer := bufio.NewReader(session)
	for {
		res, _, err := buffer.ReadRune()
		if err != nil {
			return
		}
		if res == etx || res == eot {
			return
		}
	}
}
