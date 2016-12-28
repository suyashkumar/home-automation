import React, { Component } from 'react'; 
import AppBar from 'react-toolbox/lib/app_bar'; 
import 'react-toolbox/lib/commons.scss'; 
import Test from '../containers/Test'; 

class App extends Component {
  render() {
    return (
      <div className="App">
	  	<AppBar title="ADPL Dashboard"/> 
		<Test />
		hi
      </div>
    );
  }
}

export default App;
