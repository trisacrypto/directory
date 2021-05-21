import './App.css';
import React from 'react';
import gds from './api/gds';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';
import Lookup from './components/Lookup';
import Alerts from './components/Alerts';
import Registration from './components/Registration';
import LookupResults from './components/LookupResults';


class App extends React.Component {
  state = { results: {}, alerts: [] };

  onLookup = async (query, inputType) => {
    try {
      const response = await gds.lookup(query, inputType)
      this.setState({ results: response })
    } catch(err) {
      this.setState(prevState => ({
        results: {},
        alerts: [...prevState.alerts, {variant: 'danger', message: err.message}]
      }));
      console.warn(err);
    }
  }

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

        <Tabs defaultActiveKey="directory" id="tab-nav" className="justify-content-center">
          <Tab eventKey="directory" title="Directory">
            <Row className="py-3">
              <Col md={{span: 8, offset: 2}}>
                <Lookup onSubmit={this.onLookup} />
              </Col>
            </Row>
            <Row>
              <Col>
                <LookupResults results={this.state.results} />
              </Col>
            </Row>
          </Tab>
          <Tab eventKey="register" title="Register">
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
