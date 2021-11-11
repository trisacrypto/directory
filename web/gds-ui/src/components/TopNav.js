import React from 'react';
import Container from 'react-bootstrap/Container';
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav'
import { isTestNet } from '../lib/testnet';
import { Trans } from "@lingui/macro";
import LanguageSelect from './select/LanguageSelect';
import { t } from "@lingui/macro";


const testNet = isTestNet();

const getDirectoryURL = () => {
  if (isTestNet()) {
      return [
          "https://vaspdirectory.net",
          t`Switch to Production`,
          t`You're currently on the TestNet Directory`,
      ]    
  }

  return [
          "https://trisatest.net",
          t`Switch to Production`,
          t`You're currently on the Production Directory`,
      ]
};

const [ directoryURL, directoryURLText, directoryURLTitle ] = getDirectoryURL();
const TopNav = () => {
  return (
    <Navbar variant="white" >
      <Container>
        <Navbar.Brand href="/">
          <img src="trisa-logo.jpg" alt="TRISA" height="30px" />
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="header-links" />
        <Navbar.Collapse id="header-links" className="justify-content-end">
          <Nav>
            <Nav.Link href="https://trisa.io/"><Trans>About TRISA</Trans></Nav.Link>
            <Nav.Link href="https://trisa.dev/"><Trans>Documentation</Trans></Nav.Link>
            <Nav.Link href={directoryURL} title={directoryURLTitle}>{directoryURLText}</Nav.Link>
            <LanguageSelect />
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

export default TopNav;