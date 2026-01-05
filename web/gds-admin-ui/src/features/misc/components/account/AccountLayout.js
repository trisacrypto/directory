import PropTypes from 'prop-types';
import React, { useEffect } from 'react';
import { Card, Col, Container, Row } from 'react-bootstrap';

import { getDirectoryLogo } from '@/utils';

const AccountLayout = ({ children }) => {
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
                <Card.Header
                  className="pt-4 pb-4 text-center bg-primary"
                  style={{ background: 'linear-gradient(90deg,#24a9df,#1aebb4)' }}
                >
                  <span>
                    <img src={getDirectoryLogo()} alt="" height="38" />
                  </span>
                </Card.Header>
                <Card.Body className="p-4">{children}</Card.Body>
              </Card>
            </Col>
          </Row>
        </Container>
      </div>
      <footer className="footer footer-alt">
        Created and maintained by <a href="https://rotational.io">Rotational</a> Labs in partnership
        with <a href="https://ciphertrace.com">CipherTrace</a> on behalf of{' '}
        <a href="https://travelrule.io">TRISA</a>.
      </footer>
    </>
  );
};

AccountLayout.propTypes = {
  children: PropTypes.node.isRequired,
};

export default AccountLayout;
