import React from 'react';
import './App.scss';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';
import TopNav from './components/TopNav';
import Footer from './components/Footer';
import Lookup from './components/Lookup';
import Alerts from './components/Alerts';
import Registration from './components/Registration';
import VerifyContact from './components/VerifyContact';
import Route from './components/nav/Route';
import NoRoute from './components/nav/NoRoute';
import MultiRoute from './components/nav/MultiRoute';
import { isTestNet } from './lib/testnet';

const testNet = isTestNet();
const mainRoutes = new Set(["/", "/register"]);
const allRoutes = new Set(["/", "/register", "/verify-contact"]);

// Unchanging constants that should be configured when the react app is built.
const headerClassName = testNet ? "bg-testnet-gradient" : "bg-gradient";
const headerTitle = testNet ? "TRISA TestNet Directory" : "TRISA Global Directory Service";
const headerLead = testNet ? "Get started with the TRISA TestNet to implement your Travel Rule compliance service." : "Become a TRISA certified Virtual Asset Service Provider.";
const switchLink = testNet ? <a href="https://vaspdirectory.net">Production Directory Service</a> : <a href="https://trisatest.net">TestNet Directory Service</a>

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
    console.log(testNet);
    return (
      <>
      <TopNav />
      <header className={headerClassName}>
        <div className="container">
          <div className="text-center hero">
            <h1>{headerTitle}</h1>
            <p className="lead">{headerLead}</p>
            <small>Looking for the {switchLink}?</small>
          </div>
        </div>
      </header>
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
            <Tab eventKey="/" title="Directory">
              <Lookup onAlert={this.onAlert} />
            </Tab>
            <Tab eventKey="/register" title="Register">
              <Registration onAlert={this.onAlert} />
            </Tab>
          </Tabs>
        </MultiRoute>

        <NoRoute paths={allRoutes}>
          <Row>
          <Col md={{span: 6, offset: 3}} className="text-center">
            <p className="big-number">404</p>
            <h4>PAGE NOT FOUND</h4>
            <p className="text-muted">
              The page you are looking for might have been removed, had its name changed, or is temporarily unavailable.
            </p>
            <a href="/" className="btn btn-secondary mt-2">Directory Home</a>
          </Col>
        </Row>
        </NoRoute>

      </main>
      <Footer />
      </>
    );
  }
}

export default App;
