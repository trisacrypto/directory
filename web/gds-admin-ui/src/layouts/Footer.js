import { Col, Row } from 'react-bootstrap';

import config from '../config';
import { useGetAppVersion } from '@/hooks/useGetVersion';

const Footer = () => {
    const currentYear = new Date().getFullYear();
    const { data: appVersion } = useGetAppVersion();

    return (
        <footer className="footer">
            <div className="container-fluid">
                <Row>
                    <Col md={12}>
                        {currentYear} © Created and maintained by <a href="https://rotational.io">Rotational Labs</a> in
                        partnership with <a href="https://ciphertrace.com">CipherTrace</a> on behalf of{' '}
                        <a href="https://travelrule.io">TRISA</a>.
                    </Col>
                    <Col className="d-flex">
                        {appVersion?.version ? (
                            <p data-testid="api-version">API version: {appVersion?.version}</p>
                        ) : null}
                        <>
                            {config.appVersion ? (
                                <p data-testid="app-version">&nbsp;&middot;&nbsp;App version: {config.appVersion}</p>
                            ) : null}
                            {config.gitVersion ? (
                                <p data-testid="git-version">&nbsp;&middot;&nbsp;GIT version: {config.gitVersion}</p>
                            ) : null}
                        </>
                    </Col>
                </Row>
            </div>
        </footer>
    );
};

export default Footer;
