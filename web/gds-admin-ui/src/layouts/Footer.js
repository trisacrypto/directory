import React from 'react';
import { Row, Col } from 'react-bootstrap';

const Footer = () => {
    const currentYear = new Date().getFullYear();
    return (
        <React.Fragment>
            <footer className="footer">
                <div className="container-fluid">
                    <Row>
                        <Col md={12}>{currentYear} ©
                            Created and maintained by <a href="https://rotational.io">Rotational</a> Labs in partnership with <a href="https://ciphertrace.com">CipherTrace</a> on behalf of <a href="https://trisa.io">TRISA</a>.
                        </Col>
                    </Row>
                </div>
            </footer>
        </React.Fragment>
    );
};

export default Footer;
