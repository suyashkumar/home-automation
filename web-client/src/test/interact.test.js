import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import ReactTestUtils from 'react-addons-test-utils'
import { shallow } from 'enzyme'; 
import axios from '../axios-auth';
import Interact from '../components/interact'; 
import sinon from 'sinon';

it('Renders without crashing', () => { 
  const div = document.createElement('div');
  ReactDOM.render(<Interact />, div);
});


it('Axios Call Made on button click', () => { 
	const c = ReactTestUtils.renderIntoDocument(<Interact />);
	const mockFunc = jest.fn(); 
	const spy = sinon.spy(axios, 'get');
	const button = ReactTestUtils.findRenderedDOMComponentWithClass(c, 'button') 
	ReactTestUtils.Simulate.click(button); 
	expect(spy.called);
	spy.restore(); 
});


