import React from 'react';
import Container from 'react-bootstrap/Container';
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav'
import { isTestNet } from '../lib/testnet';
import { Trans } from "@lingui/macro";
import LanguageSelect from './select/LanguageSelect';
import { i18n } from "@lingui/core";
import { t } from "@lingui/macro";


const testNet = isTestNet();
const directoryURL = testNet ? "https://vaspdirectory.net/" : "https://trisatest.net/";
const directoryURLTitle = `You're currently on the ${testNet ? "TestNet" : "Production"} Directory`;
const directoryURLText = i18n._(t`Switch to`+` ${testNet ? i18n._(t`Production`) : i18n._(t`TestNet`)}`);

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