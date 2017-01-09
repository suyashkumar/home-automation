import React, { Component } from 'react'; 
import Drawer from 'react-toolbox/lib/drawer';
import { List, ListItem, ListSubHeader, ListDivider, ListCheckbox } from 'react-toolbox/lib/list';


const NavigationMenu = props => { 
	const getNavigateAndCloseFunc = url => {
		return () => {
			window.open(url, '_self');
			props.handleToggle();
		}
	}

	return (
		<Drawer active={props.active} onOverlayClick={props.handleToggle}> 
			<div className="">
				<List selectable ripple>
					<ListSubHeader caption='Navigation' />
					<ListItem 
						caption="Home" 
						leftIcon="home" 
						className="material"
						onClick={getNavigateAndCloseFunc('/#/')} />
					<ListItem 
						caption="Interact" 
						leftIcon="settings_input_antenna" 
						className="material"
						onClick={getNavigateAndCloseFunc('/#/interact')} />
					<ListItem 
						caption="Your Account" 
						leftIcon="account_box" 
						className="material" />
					<ListItem 
						caption="Login" 
						leftIcon="trending_flat" 
						onClick={getNavigateAndCloseFunc('/#/login')} />
				</List>
			</div>
		</Drawer>
	);

}; 

export default NavigationMenu;
