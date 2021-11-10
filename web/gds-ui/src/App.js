import React from 'react';
import './App.scss';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';
import TopNav from './components/TopNav';
import Hero from './components/hero';
import Footer from './components/Footer';
import Lookup from './components/Lookup';
import Alerts from './components/Alerts';
import NotFound from './components/NotFound';
import Registration from './components/Registration';
import VerifyContact from './components/VerifyContact';
import Route from './components/nav/Route';
import NoRoute from './components/nav/NoRoute';
import MultiRoute from './components/nav/MultiRoute';
import { NetworkStore } from './contexts/NetworkContext';
import { i18n } from "@lingui/core";
import { t } from "@lingui/macro";

const mainRoutes = new Set(["/", "/register"]);
const allRoutes = new Set(["/", "/register", "/verify-contact"]);


class App extends React.Component {
  state = { alerts: [], currentPath: window.location.pathname };

  onAlert = (variant, message) => {
    this.setState(prevState => ({
      alerts: [...prevState.alerts, {variant: variant, message: message}]
    }));
  }

  onDismissAlert = (idx) => {
    let alerts = [...this.state.alerts];
    alerts.splice(idx, 1);
    this.setState({alerts: alerts});
  }

  onTabSelect = (key) => {
    window.history.pushState({}, '', key);
    this.setState({ currentPath: window.location.pathname });
  }


  render() {
    return (
      <NetworkStore>
        <TopNav />
        <Hero />
        <main role="main" className="overlap container">
          <Row>
            <Col md={{span: 8, offset: 2}}>
              <Alerts alerts={this.state.alerts} onDismiss={this.onDismissAlert} />
            </Col>
          </Row>

          <Route path="/verify-contact">
            <VerifyContact />
          </Route>

          <MultiRoute paths={mainRoutes}>
            <Tabs activeKey={this.state.currentPath} id="main-tab-nav" className="justify-content-center" onSelect={this.onTabSelect}>
              <Tab eventKey="/" title={i18n._(t`Directory`)}>
                <Lookup onAlert={this.onAlert} />
              </Tab>
              <Tab eventKey="/register" title={i18n._(t`Register`)}>
                <Registration onAlert={this.onAlert} />
              </Tab>
            </Tabs>
          </MultiRoute>

          <NoRoute paths={allRoutes}>
            <NotFound />
          </NoRoute>

        </main>
        <Footer />
      </NetworkStore>
    );
  }
}

export default App;
