import React from 'react';

const LookupResults = (props) => {
  let results = "";
  if (!(props.results && Object.keys(props.results).length === 0 && props.results.constructor === Object)) {
    results = JSON.stringify(props.results, null, 2);
  }
  return (
    <div className="lookup-results">
      <pre>{ results }</pre>
    </div>
  );
}

export default LookupResults;