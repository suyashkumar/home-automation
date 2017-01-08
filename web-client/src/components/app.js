import React, { Component } from 'react'; 
import AppBar from 'react-toolbox/lib/app_bar'; 
import 'react-toolbox/lib/commons.scss'; 
import Test from '../containers/Test'; 
import NavigationMenu from './app/navigation-menu';

class App extends Component { 

	state = {
		navMenuOpen: false,
	};

	navMenuToggle = () => {
		this.setState({navMenuOpen: !this.state.navMenuOpen});
	}

  	render() {
    	return (
      	<div className="App">
	  		<AppBar 
				title="Conduit" 
				leftIcon="menu" 
				onLeftIconClick={this.navMenuToggle}
				className="material"/> 

			<NavigationMenu
				active={this.state.navMenuOpen}
				handleToggle={this.navMenuToggle} /> 
			{this.props.children}
      	</div>
    	);
  	}
}

export default App;
