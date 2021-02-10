import React, { useState } from 'react';
import logo from './logo.svg';
import './App.css';


interface response {
  time: string
}
function App() {
  var [time, setTime] = useState<string>("")

  var getTime = () => {

    fetch('/api')
      .then(r => {
        if (!r.ok) return Promise.reject()
        return r.json()
      })
      .then(j => {
        let r: response = j
        setTime(r.time)
      })
      .catch(e => console.error(e))

  }

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <span>{time}</span>
        <button onClick={getTime}>Get Time</button>
      </header>
    </div>
  );
}

export default App;
