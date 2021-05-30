import './App.css';
import React from 'react';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';
import Lookup from './components/Lookup';
import Alerts from './components/Alerts';
import Registration from './components/Registration';


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
      <>
      <main role="main" className="container">
        <header className="pt-5 pb-3 text-center">
          <img className="d-block mx-auto mb-4" src="logo192.png" alt="" width="72" height="72" />
          <h2>TRISA Directory Service</h2>
          <p className="lead">Lookup Virtual Asset Service Providers that are TRISA certified.</p>
        </header>

        <Row>
          <Col md={{span: 8, offset: 2}}>
            <Alerts alerts={this.state.alerts} onDismiss={this.onDismissAlert} />
          </Col>
        </Row>

        <Tabs activeKey={this.state.currentPath} id="main-tab-nav" className="justify-content-center" onSelect={this.onTabSelect}>
          <Tab eventKey="/" title="Directory">
            <Lookup onAlert={this.onAlert} />
          </Tab>
          <Tab eventKey="/register" title="Register">
            <Registration onAlert={this.onAlert} />
          </Tab>
        </Tabs>

      </main>
      <footer className="footer">
        <div className="container text-center">
          <span className="text-muted">A demonstration of the <a href="https://trisa.io/">TRISA</a> architecture for Cryptocurrency Travel Rule compliance.</span>
        </div>
      </footer>
      </>
    );
  }
}

export default App;
