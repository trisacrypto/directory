// @flow
import React from 'react';
import { Row, Col } from 'react-bootstrap';

const Footer = (): React$Element<any> => {
    const currentYear = new Date().getFullYear();
    return (
        <React.Fragment>
            <footer className="footer">
                <div className="container-fluid">
                    <Row>
                        <Col md={12}>{currentYear} ©
                            Created and maintained by <a href="https://www.rotatinal.io">Rotational</a> Labs in partnership with <a href="https://www.cyphertrace.com">CipherTrace</a> on behalf of <a href="https://www.trisa.io">TRISA</a>.
                        </Col>
                    </Row>
                </div>
            </footer>
        </React.Fragment>
    );
};

export default Footer;
