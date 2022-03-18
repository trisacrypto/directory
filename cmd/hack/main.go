/*
Hack to try to reissue certificates and correctly update alice and bob.
*/
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/gds/store"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

var (
	db   store.Store
	conf config.Config
)

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	var err error

	// Load the configuration from the environment
	if conf, err = config.New(); err != nil {
		log.Fatal(err)
	}
	conf.Database.ReindexOnBoot = false

	// Connect to the trtl server and create a store to access data directly like GDS
	if db, err = store.Open(conf.Database); err != nil {
		log.Fatal(err)
	}

	// JOB 1: update the Alice and Bob records if we're on the TestNet.
	if err = updateRobotCerts(); err != nil {
		log.Fatal(err)
	}

	// JOB 2: reissue certificates for NMS
	if err = reissueNMSCerts(); err != nil {
		log.Fatal(err)
	}
}

func updateRobotCerts() (err error) {
	names := map[string]string{
		"alice": "7a96ca2c-2818-4106-932e-1bcfd743b04c",
		"bob":   "9e069e01-8515-4d57-b9a5-e249f7ab4fca",
	}

	for name, vaspID := range names {
		// Lookup the VASP record
		var vasp *pb.VASP
		if vasp, err = db.RetrieveVASP(vaspID); err != nil {
			return fmt.Errorf("unable to retrieve %s record: %v", name, err)
		}

		// Unmarshal the identity certificate
		var data []byte
		if data, err = base64.StdEncoding.DecodeString(certs[name]); err != nil {
			return fmt.Errorf("could not deocde base64 for %s: %v", name, err)
		}

		vasp.IdentityCertificate = new(pb.Certificate)
		if err = proto.Unmarshal(data, vasp.IdentityCertificate); err != nil {
			return fmt.Errorf("could not unmarshal %s certificate: %v", name, err)
		}

		// Update the VASP status as verified/certificate issued
		if err = models.UpdateVerificationStatus(vasp, pb.VerificationState_VERIFIED, "manual certificate reissuance", "benjamin@rotational.io"); err != nil {
			return fmt.Errorf("could not update VASP verification status %s: %v", name, err)
		}
		if err = db.UpdateVASP(vasp); err != nil {
			return fmt.Errorf("could not update VASP %s: %v", name, err)
		}

		// Create certificate request for the issued certs
		if err = db.UpdateCertReq(certreqs[name]); err != nil {
			return fmt.Errorf("could not save certificate request for VASP %s: %v", name, err)
		}
	}
	return nil
}

var certs = map[string]string{
	"alice": "CAMSEGC53MBF1ZYImUrSF5/EpEQagAREgXemG6RJ3ubLvYAWNpGZIPCEvMaUseoboRQOcFkze0teJTyJsOgl0tbR5R3jerVSn1UfZzrVmj4A2rGSS/SJa/9x1/YzSO7J1utTXU+/lDUcWGWOZ3X/QOFomifWg3KfgDuB2j4dgALWtwdSidOTXrm8Sxsj0fJTGpiG6Usr8Hkft7FkFxa8hhG0GHbxPljW/GZ6u8lTtUCaSWmmM6rA3ixiI55wxxH4xC9EknM9PTzL1yTumx58Xf6ogCsMPPJHH5AcSuK3wTl0809Lx4pR0o8w0SX0E8t/pqm4/WyKxSnVCjn2v0q/ThM6lYlY0q2f5MwXtG3yGOzH2duBwjpQ3L8N5/n+Jf+GWoE7Rh/ClPKlEh/hHjtnwD/9CXykHMFMXLHQAaziE9n8UaYQ9X77vqnE/QYwjtMfK9OY+NV6lzhq6tGk9Ed2oqVdJQtGvZOO9QfGzQWKzNpfncbSXS+eLrjtzQB7VOhxJyfGajBRg93xjvd8ZDscS2I0SVriZ0UxGzpfJEEZmpYu6BdM0637gr+TBYPAPayFo4iDYXxfravpO6PTPZxQAwOtoouvw7R+p5/06jtzIqYVCu1vJAf3L2XIdIo1Kbdxjlqfv1fGGe9TlA1CCY0l/IsMVTKUElO9J5k/wDiACdNEIyXNclbccDhiWDIBESKSTAp7YRQxsiIKU0hBMzg0LVJTQSoDUlNBMkQKFWFwaS5hbGljZS52YXNwYm90Lm5ldBoPQ2lwaGVyVHJhY2UgSW5jMgpNZW5sbyBQYXJrOgpDYWxpZm9ybmlhSgJVUzpFChZDaXBoZXJUcmFjZSBJc3N1aW5nIENBGg9DaXBoZXJUcmFjZSBJbmMyCk1lbmxvIFBhcms6CkNhbGlmb3JuaWFKAlVTQhQyMDIyLTAyLTI4VDE3OjM5OjM4WkoUMjAyMy0wMy0yOFQxNzozOTozN1pawBMtLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTUlJSEJqQ0NCTzZnQXdJQkFnSVFZTG5jd0VYVmxnaVpTdElYbjhTa1JEQU5CZ2txaGtpRzl3MEJBUXdGQURCeQpNUXN3Q1FZRFZRUUdFd0pWVXpFVE1CRUdBMVVFQ0JNS1EyRnNhV1p2Y201cFlURVRNQkVHQTFVRUJ4TUtUV1Z1CmJHOGdVR0Z5YXpFWU1CWUdBMVVFQ2hNUFEybHdhR1Z5VkhKaFkyVWdTVzVqTVI4d0hRWURWUVFERXhaRGFYQm8KWlhKVWNtRmpaU0JKYzNOMWFXNW5JRU5CTUI0WERUSXlNREl5T0RFM016a3pPRm9YRFRJek1ETXlPREUzTXpregpOMW93Y1RFTE1Ba0dBMVVFQmhNQ1ZWTXhFekFSQmdOVkJBZ1RDa05oYkdsbWIzSnVhV0V4RXpBUkJnTlZCQWNUCkNrMWxibXh2SUZCaGNtc3hHREFXQmdOVkJBb1REME5wY0dobGNsUnlZV05sSUVsdVl6RWVNQndHQTFVRUF4TVYKWVhCcExtRnNhV05sTG5aaGMzQmliM1F1Ym1WME1JSUNJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBZzhBTUlJQwpDZ0tDQWdFQXhWZVpjdGJ5YzhWL3MwUW5MZXRpOEJ1ekJiMWR0cDhMTzVETHBzUHd6dHR5VWNIZWNpZ0NjNjNKCk9wR0VpSVh4dHI0T1BRN2Z5bHNkd3dTQ1ZDUEtRNVByZGJHYzlwQW1LTUk0Nm9TTjN2OWQ3dFRXSm5wMFladlcKUytadXJYNVRoL0RyakJGZ2ZTdnhkVVZrQTVpcjJtcStrUUt6bkQvTlFXbU9ETXhqdFh6MHFWaWQrQlk0S2c3awpEd2JHY0VUL3FoR0lJTFh1N2NGK0oya3FUenFEM2pnaEVkOGxOMllTdm5ySXdYbkNaRjJEZVRINDVBTmtLeHpTCmw3d0owY2lBU1pvZUZKUDRNdjNFdEJ5UHBKWnRuMlZlWlJETktGNzdaMTVMYkdTOFl3QUJMZUNyUzhIV1ZGVk0KMGFBUlNkWjRpYkJudWxockdUdkJTdmpVQWtmbXpsSmVKSitNQWxvQ0tRazFKVy8xTElReTNtczhHNlVXZWJqVQpZelBxT3RDdTN5TFJVeDNmQzM4ZzNvRnFibmlxdUYySVpQWThueGVBak1COS9iWGZWVlltWXpPSHhWQW0yMG5MCm9LR0ZrVkx0NWptYlNTb0xOMGtZNEtuN1Q5dFZ4TzMzcTlpaGtQTnl4SlIzeDBwQzFoT05LU1JUb2dXeEpsSXQKM01XYkUxZCtGMGlPd0NjMTk2ZlpOdG9seVlpTW1jK05Sd2lpdC8vc0Q4R3dabkU5UzRpQzcwTTc3bTdzNGE1cQpjN2tKNHlIdnJSb3p2SEhSQnFqdEdMbXArQ0dNVW9CUVk2SVJhMUxMeVRWbjA5SmVMdGoxejlhNDFyNVoyeHppCmJKcHZiQnN2Vzh3Q2U2cHpVb2JOd2RBbjhBaEt1NTAzV1laSjNqRWVOUW8rVWJ3Q01lOENBd0VBQWFPQ0FaY3cKZ2dHVE1COEdBMVVkSXdRWU1CYUFGRnU3aU03WkdEbkxpMnVDOFFCdldtUW00LzEyTUE0R0ExVWREd0VCL3dRRQpBd0lGNERBTUJnTlZIUk1CQWY4RUFqQUFNQ0FHQTFVZEpRRUIvd1FXTUJRR0NDc0dBUVVGQndNQ0JnZ3JCZ0VGCkJRY0RBVEJWQmdOVkhSOEVUakJNTUVxZ1NLQkdoa1JvZEhSd09pOHZZMmx3YUdWeWRISmhZMlV1WTNKc0xtbHYKZEM1elpXTjBhV2R2TG1OdmJTOURhWEJvWlhKVWNtRmpaVWx1ZEdWeWJXVmthV0YwWlVOQkxtTnliRENCbHdZSQpLd1lCQlFVSEFRRUVnWW93Z1ljd1VBWUlLd1lCQlFVSE1BS0dSR2gwZEhBNkx5OWphWEJvWlhKMGNtRmpaUzVqCmNuUXVhVzkwTG5ObFkzUnBaMjh1WTI5dEwwTnBjR2hsY2xSeVlXTmxTVzUwWlhKdFpXUnBZWFJsUTBFdVkzSjAKTURNR0NDc0dBUVVGQnpBQmhpZG9kSFJ3T2k4dlkybHdhR1Z5ZEhKaFkyVXViMk56Y0M1cGIzUXVjMlZqZEdsbgpieTVqYjIwd0lBWURWUjBSQkJrd0Y0SVZZWEJwTG1Gc2FXTmxMblpoYzNCaWIzUXVibVYwTUIwR0ExVWREZ1FXCkJCUnovb2s3TW5yM25Pbzh6N0xhTmtEWlc5L2pzekFOQmdrcWhraUc5dzBCQVF3RkFBT0NBZ0VBUklGM3BodWsKU2Q3bXk3MkFGamFSbVNEd2hMekdsTEhxRzZFVURuQlpNM3RMWGlVOGliRG9KZExXMGVVZDQzcTFVcDlWSDJjNgoxWm8rQU5xeGtrdjBpV3YvY2RmMk0wanV5ZGJyVTExUHY1UTFIRmhsam1kMS8wRGhhSm9uMW9OeW40QTdnZG8rCkhZQUMxcmNIVW9uVGsxNjV2RXNiSTlIeVV4cVlodWxMSy9CNUg3ZXhaQmNXdklZUnRCaDI4VDVZMXZ4bWVydkoKVTdWQW1rbHBwak9xd040c1lpT2VjTWNSK01RdlJKSnpQVDA4eTljazdwc2VmRjMrcUlBckREenlSeCtRSEVyaQp0OEU1ZFBOUFM4ZUtVZEtQTU5FbDlCUExmNmFwdVAxc2lzVXAxUW81OXI5S3YwNFRPcFdKV05LdG4rVE1GN1J0CjhoanN4OW5iZ2NJNlVOeS9EZWY1L2lYL2hscUJPMFlmd3BUeXBSSWY0UjQ3WjhBLy9RbDhwQnpCVEZ5eDBBR3MKNGhQWi9GR21FUFYrKzc2cHhQMEdNSTdUSHl2VG1QalZlcGM0YXVyUnBQUkhkcUtsWFNVTFJyMlRqdlVIeHMwRgppc3phWDUzRzBsMHZuaTY0N2MwQWUxVG9jU2NueG1vd1VZUGQ4WTczZkdRN0hFdGlORWxhNG1kRk1SczZYeVJCCkdacVdMdWdYVE5PdCs0Sy9rd1dEd0Qyc2hhT0lnMkY4WDYycjZUdWowejJjVUFNRHJhS0xyOE8wZnFlZjlPbzcKY3lLbUZRcnRieVFIOXk5bHlIU0tOU20zY1k1YW43OVh4aG52VTVRTlFnbU5KZnlMREZVeWxCSlR2U2VaUDhBNApnQW5UUkNNbHpYSlczSEE0WWxneUFSRWlra3dLZTJFVU1iST0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQpizyIfiwgAAAAAAAD/zJjH8utKkt73eArtGQrCm8UsCt4T3u3gCA8QIPzTK/7n3L59e6ZnRlJHKMRllkl+hagvf5X/8+fHCpJi/g9OcDxFVDjgCb+ikKEoMttyHPsiK3AoLKgUO9bH/BCioK+axF2VaKTdzuGByVbdXHeNxBwwC+xDBDx7QYb9PTg75gPbloRDDfxb8AxWkADiCxxraDYqftMw2fOB+MTe38fY09C8MNigTKIrXxKv9BZig41/L6wNy0b7I5WCK5DVOkb9yg2J1nDoQ/6djBfOhE8jdoKSSPXzQWwTl1VzzETSkBgVwWQNFo94T7kMXrlevIAZd3e/xOkndhu88WcMMpHpyD1BN0D3+6/VBhcExincwGErM2BB5XGdWWdSP2SYuqWh8Jex3IO4Dumz4dwVka3z4XtKPAh/D04eD5ufXKr7vHeuODR7Rei3+BZKgz1+JQOnEUBxxH704eegzF4fkzrH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",
	"bob":   "CAMSEB5ZHHcuDf+zq9ac0Vm+ujgagARv1TdUlFIF26KaLYi8LNLe8Z2MDe6EBiK1URXDyShqwV3QIBGibBdg6+kP/OMC1AxebnqXMQWwAnHoKgwJ1vXt58Ak4d3EWcHrYlK3KsGKQ+catN2WuWataXwzZoJKBlD8KwST3FFvY+0mgCEg2FK2fkX1xDe4yJrG3UEQVA+EHcfJrFWwujuR+WlT1uXedcpOhv23nHTlobxwr+y4yiJBlefXfROZQ7FA7iBpHuaDKmZAXv2vHnGF0Pk9L8dLfSOB34kMd2MJEp6I5M0oHBr0EL1LOZ6hP2HwPdvXsxwcFi/6VIMXyz9aJ1P25mKLNRX41ebA89YAc83v/V7ubbPbdlXkj4OZRudmtpA0dK9VzErdX6H6N02c/1m9Xl3Rcr4y1PdWKL2NxbK1K8MtkILRasHHYQVvcr+Bj+lrNvydOBeeVkjK025vAnGz2308rZ9ICtdEmNmqFBMYyYmO+GsYdq4x1Ak8nKybqVbS7NZsQAq+f0g5thSI1vRSI3nopGX9wmlLO9taybZGK6CxID57KB+mJZYCFkkN6gok614cUzEW2/ISycqk+YCJXry/4eqpq+bIHX8O7Bpk93WDcmvhZbJYtpjJoHHXbWS4AqsQnyE6yMcRMO0s1OuDyZpySLyJ3oW5cAVT0oDcQfwHSpdZElzoJHjS7dZXujdIN6h1QSIKU0hBMzg0LVJTQSoDUlNBMkIKE2FwaS5ib2IudmFzcGJvdC5uZXQaD0NpcGhlclRyYWNlIEluYzIKTWVubG8gUGFyazoKQ2FsaWZvcm5pYUoCVVM6RQoWQ2lwaGVyVHJhY2UgSXNzdWluZyBDQRoPQ2lwaGVyVHJhY2UgSW5jMgpNZW5sbyBQYXJrOgpDYWxpZm9ybmlhSgJVU0IUMjAyMi0wMi0yOFQxNzo1MzozNVpKFDIwMjMtMDMtMjhUMTc6NTM6MzRaWrwTLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUhBakNDQk9xZ0F3SUJBZ0lRSGxrY2R5NE4vN09yMXB6UldiNjZPREFOQmdrcWhraUc5dzBCQVF3RkFEQnkKTVFzd0NRWURWUVFHRXdKVlV6RVRNQkVHQTFVRUNCTUtRMkZzYVdadmNtNXBZVEVUTUJFR0ExVUVCeE1LVFdWdQpiRzhnVUdGeWF6RVlNQllHQTFVRUNoTVBRMmx3YUdWeVZISmhZMlVnU1c1ak1SOHdIUVlEVlFRREV4WkRhWEJvClpYSlVjbUZqWlNCSmMzTjFhVzVuSUVOQk1CNFhEVEl5TURJeU9ERTNOVE16TlZvWERUSXpNRE15T0RFM05UTXoKTkZvd2J6RUxNQWtHQTFVRUJoTUNWVk14RXpBUkJnTlZCQWdUQ2tOaGJHbG1iM0p1YVdFeEV6QVJCZ05WQkFjVApDazFsYm14dklGQmhjbXN4R0RBV0JnTlZCQW9URDBOcGNHaGxjbFJ5WVdObElFbHVZekVjTUJvR0ExVUVBeE1UCllYQnBMbUp2WWk1MllYTndZbTkwTG01bGREQ0NBaUl3RFFZSktvWklodmNOQVFFQkJRQURnZ0lQQURDQ0Fnb0MKZ2dJQkFOeHBqVGRoMisweXlMakRFYUxlRUZEYXJaUmg0TExPNThNbHFseFBVTCtwcTdrano1QVpSelRjb1NlZgpxQ0V1aHpDWi9UdlZOcUV3ZWt1Tk45T01jc3NmSk5lYW1jN3RYSGltbFptWUVwRVZ3eW10NjREQ2orUFlIWDFMCmRyRldEZ09Kc1BGMDI4ZFFuVWljMXFWVy9DOWNLSnBSa25lZGNJbUZZanJDNWNnQkZZZnQvRnZ2bVVGL3Q0cEcKTWJHYWoyZklKWUVYdGt6b1p2dElvNjg0Qzd5L2t2NzJVdmpiU3RpcUZzVkIvaUhrSVd6ZzJwdTJiQ1VFR0RaagpoNTV6aDhITXloMEwzeEZ2RGxuUXl1eWdMYUlNRU1EVHlhTksrRlpZVHBNcWNXVjdVR0o1ZGFiZ2NiLzIvRUhPClZoKzJ2dkNpc2QzcTV5YmlyQzZlUVJ5UXFFYnFkNzd3NHRaRFVyOThvUytEemdHZ2VVdWY3OVpJZ1dYOXg1S0UKRkFaOWpsdzRtRnVPQVE5dTYzMU5iT0llYVFZT3RLZ05RNGg1OVhHQzBjRHFxWFc1NlFkeFAweVFDb3hNODYwawplVjVNNkdoV3E2QUJudVdJN1dieVBTaHhvaWhyUzd1cHVtZE11RHBTRHQwYWxsR0ZmbTFLeXNZUkZiL1N6NXFCCjRYMkUrSk54T0daUUxsT2FpQm1icXdPeE9wY2lndEVjVlc1QkVSVDJ6aVMrdGttUEVqRjFxeHMySFVCa3lHL04KYWhvNjBFZVFzYWgvZWNlVzJuMVllWlFPYlNpcWU2dkd0TXJPeFY3ZmNGc24vckoxdVBzSlpnL0l4aWxOSm5jbQpnUWk1S29IK080NVJYUU5jM0lUejdDVmdmd0VyaE1kUW9PZy91RmJ0dytkaHlIZlhBZ01CQUFHamdnR1ZNSUlCCmtUQWZCZ05WSFNNRUdEQVdnQlJidTRqTzJSZzV5NHRyZ3ZFQWIxcGtKdVA5ZGpBT0JnTlZIUThCQWY4RUJBTUMKQmVBd0RBWURWUjBUQVFIL0JBSXdBREFnQmdOVkhTVUJBZjhFRmpBVUJnZ3JCZ0VGQlFjREFnWUlLd1lCQlFVSApBd0V3VlFZRFZSMGZCRTR3VERCS29FaWdSb1pFYUhSMGNEb3ZMMk5wY0dobGNuUnlZV05sTG1OeWJDNXBiM1F1CmMyVmpkR2xuYnk1amIyMHZRMmx3YUdWeVZISmhZMlZKYm5SbGNtMWxaR2xoZEdWRFFTNWpjbXd3Z1pjR0NDc0cKQVFVRkJ3RUJCSUdLTUlHSE1GQUdDQ3NHQVFVRkJ6QUNoa1JvZEhSd09pOHZZMmx3YUdWeWRISmhZMlV1WTNKMApMbWx2ZEM1elpXTjBhV2R2TG1OdmJTOURhWEJvWlhKVWNtRmpaVWx1ZEdWeWJXVmthV0YwWlVOQkxtTnlkREF6CkJnZ3JCZ0VGQlFjd0FZWW5hSFIwY0RvdkwyTnBjR2hsY25SeVlXTmxMbTlqYzNBdWFXOTBMbk5sWTNScFoyOHUKWTI5dE1CNEdBMVVkRVFRWE1CV0NFMkZ3YVM1aWIySXVkbUZ6Y0dKdmRDNXVaWFF3SFFZRFZSME9CQllFRk83ZQo3UExTZUdPUXI2aVhvN3RWN2xJTFR3YUxNQTBHQ1NxR1NJYjNEUUVCREFVQUE0SUNBUUJ2MVRkVWxGSUYyNkthCkxZaThMTkxlOFoyTURlNkVCaUsxVVJYRHlTaHF3VjNRSUJHaWJCZGc2K2tQL09NQzFBeGVibnFYTVFXd0FuSG8KS2d3SjF2WHQ1OEFrNGQzRVdjSHJZbEszS3NHS1ErY2F0TjJXdVdhdGFYd3pab0pLQmxEOEt3U1QzRkZ2WSswbQpnQ0VnMkZLMmZrWDF4RGU0eUpyRzNVRVFWQStFSGNmSnJGV3d1anVSK1dsVDF1WGVkY3BPaHYyM25IVGxvYnh3CnIreTR5aUpCbGVmWGZST1pRN0ZBN2lCcEh1YURLbVpBWHYydkhuR0YwUGs5TDhkTGZTT0IzNGtNZDJNSkVwNkkKNU0wb0hCcjBFTDFMT1o2aFAySHdQZHZYc3h3Y0ZpLzZWSU1YeXo5YUoxUDI1bUtMTlJYNDFlYkE4OVlBYzgzdgovVjd1YmJQYmRsWGtqNE9aUnVkbXRwQTBkSzlWekVyZFg2SDZOMDJjLzFtOVhsM1JjcjR5MVBkV0tMMk54YksxCks4TXRrSUxSYXNISFlRVnZjcitCaitsck52eWRPQmVlVmtqSzAyNXZBbkd6MjMwOHJaOUlDdGRFbU5tcUZCTVkKeVltTytHc1lkcTR4MUFrOG5LeWJxVmJTN05ac1FBcStmMGc1dGhTSTF2UlNJM25vcEdYOXdtbExPOXRheWJaRwpLNkN4SUQ1N0tCK21KWllDRmtrTjZnb2s2MTRjVXpFVzIvSVN5Y3FrK1lDSlhyeS80ZXFwcStiSUhYOE83QnBrCjkzV0RjbXZoWmJKWXRwakpvSEhYYldTNEFxc1FueUU2eU1jUk1PMHMxT3VEeVpweVNMeUozb1c1Y0FWVDBvRGMKUWZ3SFNwZFpFbHpvSkhqUzdkWlh1amRJTjZoMVFRPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQpisyIfiwgAAAAAAAD/1NjHzqxasiDgOU/R81Qr8WZwB3jvPTO8TxJvnr71733r1Dl1napLaqlzuAJWEEgR60v+98+P4UXZ+F8s73iyILO0x/9aBXRZluiOZRlzrulTZuhatqWhz4sbNd6EuUDfxwkzHDc52mDqfm76VqROkKHtU6A55gZ0ez1ZO+YC2xb5Uwn8h/d0hhdpyOdZRldtWFjTMDnyEfvG3t9jzKWrXhjsQCaStS8Kd/rwsc7Ev29sdMuGhzMVgzuQlCaG/doNsU53yFP6nYzjr4RLI2YCkkjx81HoEpdRcsSA0hD7yLzB6AwacZ5865x8mxyPGJ7+GMH0s/bonP7HGmAI05k9vKbT/e9Ha3Q2CPSLf2iHqY2AoWuP7Y0mE4cxQ5Q9Dfk/xXIPYHtoyMbrkAWmycf1Ejk6/B2cPA40vrnYDPng3HFoDDI/7PHD5zoz/UpGX7oHxBHz1UbliFsMjiPjjEcK1EZsKDiW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",
}

var certreqs = map[string]*models.CertificateRequest{
	"alice": {
		Id:           "e2d7ef62-28e2-4fb9-9fd6-4fb107d942a6",
		Vasp:         "7a96ca2c-2818-4106-932e-1bcfd743b04c",
		CommonName:   "api.alice.vaspbot.net",
		Status:       models.CertificateRequestState_COMPLETED,
		AuthorityId:  23,
		BatchId:      309867,
		BatchName:    "trisatest.net-certreq-e2d7ef62-28e2-4fb9-9fd6-4fb107d942a6",
		BatchStatus:  "READY_FOR_DOWNLOAD",
		OrderNumber:  309867,
		CreationDate: "2022-02-28T11:39:38.000+0000",
		Profile:      "CipherTrace EE",
		RejectReason: "",
		Params:       map[string]string{},
		Created:      time.Now().Format(time.RFC3339Nano),
	},
	"bob": {
		Id:           "1177a95b-7e24-411a-a0a8-3210e4b93b2b",
		Vasp:         "9e069e01-8515-4d57-b9a5-e249f7ab4fca",
		CommonName:   "api.bob.vaspbot.net",
		Status:       models.CertificateRequestState_COMPLETED,
		AuthorityId:  23,
		BatchId:      309870,
		BatchName:    "trisatest.net-certreq-1177a95b-7e24-411a-a0a8-3210e4b93b2b",
		BatchStatus:  "READY_FOR_DOWNLOAD",
		OrderNumber:  309870,
		CreationDate: "2022-02-28T11:53:42.000+0000",
		Profile:      "CipherTrace EE",
		RejectReason: "",
		Params:       map[string]string{},
		Created:      time.Now().Format(time.RFC3339Nano),
	},
}

func reissueNMSCerts() (err error) {
	var vasp *pb.VASP
	vaspID := "00634f29-e22e-48f5-be2a-74feeee33464"
	commonName := "eczesgsg.trisa.test-travelrule.sygna.io"
	endpoint := "eczesgsg.trisa.test-travelrule.sygna.io:443"

	if vasp, err = db.RetrieveVASP(vaspID); err != nil {
		return fmt.Errorf("could not retrieve vasp record: %v", err)
	}

	// Reset the Travel Rule information
	vasp.CommonName = commonName
	vasp.TrisaEndpoint = endpoint
	vasp.IdentityCertificate = nil

	// Cannot validate
	// if err = vasp.Validate(false); err != nil {
	// 	return fmt.Errorf("vasp is not valid: %v", err)
	// }

	if err = models.UpdateVerificationStatus(vasp, pb.VerificationState_REVIEWED, "manually reissuing testnet certificates", "benjamin@rotational.io"); err != nil {
		return fmt.Errorf("could not update VASP record: %v", err)
	}

	// Create a new certificate request for the reissuance
	password := secrets.CreateToken(16)
	fmt.Printf("password is %q\n", password)

	var certreq *models.CertificateRequest
	if certreq, err = models.NewCertificateRequest(vasp); err != nil {
		return err
	}

	if err = models.UpdateCertificateRequestStatus(certreq, models.CertificateRequestState_READY_TO_SUBMIT, "reissue certificates", "benjamin@rotational.io"); err != nil {
		return fmt.Errorf("could not update certreq: %v", err)
	}

	// Connect to secret manager
	var sm *secrets.SecretManager
	if sm, err = secrets.New(conf.Secrets); err != nil {
		return err
	}

	// Make a new secret of type password
	secretType := "password"
	if err = sm.With(certreq.Id).CreateSecret(context.TODO(), secretType); err != nil {
		return err
	}
	if err = sm.With(certreq.Id).AddSecretVersion(context.TODO(), secretType, []byte(password)); err != nil {
		return err
	}

	// Update the certreq
	if err = db.UpdateCertReq(certreq); err != nil {
		return err
	}

	// Append the certificate request to the VASP
	if err = models.AppendCertReqID(vasp, certreq.Id); err != nil {
		return err
	}

	// Update the vasp
	if err = db.UpdateVASP(vasp); err != nil {
		return err
	}

	return nil
}
