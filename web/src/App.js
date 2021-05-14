import './App.css';
import React from 'react';
import gds from './api/gds';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Lookup from './components/Lookup';
import LookupResults from './components/LookupResults';


class App extends React.Component {
  state = { results: {} };

  onLookup = async (query, inputType) => {
    try {
      const response = await gds.lookup(query, inputType)
      this.setState({ results: response })
    } catch(err) {
      this.setState({results: {}})
      console.log(err);
    }
  }

  render() {
    return (
      <>
      <main role="main" className="container">
        <header className="pt-5 pb-3 text-center">
          <div id="alerts"></div>
          <img className="d-block mx-auto mb-4" src="logo192.png" alt="" width="72" height="72" />
          <h2>TRISA Directory Service</h2>
          <p className="lead">Lookup Virtual Asset Service Providers that are TRISA certified.</p>
        </header>

        <Row>
          <Col md={{span: 8, offset: 2}}>
            <Lookup onSubmit={this.onLookup} />
          </Col>
        </Row>
        <Row>
          <Col>
            <LookupResults results={this.state.results} />
          </Col>
        </Row>

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
