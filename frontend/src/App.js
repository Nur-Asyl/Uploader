import logo from './logo.svg';
import './App.css';
import Uploaders from "./components/uploaders.components";
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Uploader <code>json</code>
        </p>
        <div>
        <Uploaders/>
        </div>
      </header>
    </div>
  );
}

export default App;
