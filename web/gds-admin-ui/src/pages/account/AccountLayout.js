// @flow
import React, { useEffect } from 'react';
import { Container, Row, Col, Card } from 'react-bootstrap';

import Logo from '../../assets/images/gds-trisatest-logo.png';

type AccountLayoutProps = {
    bottomLinks?: React$Element<any>,
    children?: any,
};

const AccountLayout = ({ bottomLinks, children }: AccountLayoutProps): React$Element<any> => {

    useEffect(() => {
        if (document.body) document.body.classList.add('authentication-bg');

        return () => {
            if (document.body) document.body.classList.remove('authentication-bg');
        };
    }, []);

    return (
        <>
            <div className="account-pages pt-2 pt-sm-5 pb-4 pb-sm-5">
                <Container>
                    <Row className="justify-content-center">
                        <Col md={8} lg={6} xl={5} xxl={4}>
                            <Card>
                                <Card.Header className="pt-4 pb-4 text-center bg-primary" style={{ background: "linear-gradient(90deg,#24a9df,#1aebb4)" }}>
                                    <span>
                                        <img src={Logo} alt="" height="38" />
                                    </span>
                                </Card.Header>
                                <Card.Body className="p-4">{children}</Card.Body>
                            </Card>

                            {bottomLinks}
                        </Col>
                    </Row>
                </Container>
            </div>
            <footer className="footer footer-alt">
                Created and maintained by <a href="https://rotational.io">Rotational</a> Labs in partnership with <a href="https://ciphertrace.com">CipherTrace</a> on behalf of <a href="https://trisa.io">TRISA</a>.
            </footer>
        </>
    );
};

export default AccountLayout;
