import config from '../config';
import React, { useEffect, useState } from 'react';
import { Row, Col } from 'react-bootstrap';
import getAppVersion from 'services/version';

const Footer = () => {
    const [appVersion, setAppVersion] = useState()
    const currentYear = new Date().getFullYear();

    useEffect(() => {
        getAppVersion().then(res => {
            setAppVersion(res.data)
            // eslint-disable-next-line no-console
            console.log('version', res.data.version)
        }).catch(err => {
            console.error('[getAppVersion]', err)
        })
    }, [])

    return (
        <React.Fragment>
            <footer className="footer">
                <div className="container-fluid">
                    <Row>
                        <Col md={12}>{currentYear} ©
                            Created and maintained by <a href="https://rotational.io">Rotational Labs</a> in partnership with <a href="https://ciphertrace.com">CipherTrace</a> on behalf of <a href="https://trisa.io">TRISA</a>.
                        </Col>
                        <Col className='d-flex'>
                            {appVersion?.version ? <p data-testid="api-version">API version: {appVersion?.version}</p> : null}
                            {
                                <>
                                    {config.appVersion ? <p data-testid="app-version">&nbsp;&middot;&nbsp;App version: {config.appVersion}</p> : null}
                                    {config.gitVersion ? <p data-testid="git-version">&nbsp;&middot;&nbsp;GIT version: {config.gitVersion}</p> : null}
                                </>
                            }

                        </Col>
                    </Row>
                </div>
            </footer>
        </React.Fragment>
    );
};

export default Footer;
