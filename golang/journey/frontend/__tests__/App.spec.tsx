import React from 'react';
import { render, act, waitFor, cleanup } from '@testing-library/react';
import { shallow } from 'enzyme';
import App from '../App';

let mockReturnData: { id: string; data: { x: number; y: number }[] }[];
global.fetch = jest.fn().mockResolvedValue({
  headers: { get: jest.fn(() => '0') },
  json: async () => ({ journeys: mockReturnData }),
});

describe('App', () => {
  beforeEach(() => {
    jest.useFakeTimers();
    mockReturnData = [];
  });

  it('should render with default state', () => {
    const wrapper = shallow(<App />);

    expect(wrapper.find('ResponsiveLine').prop('data')).toStrictEqual([]);
  });

  it('should render chart with data from state', () => {
    const expectedData = [
      { id: '23->42', data: [{ x: 1, y: 2 }] },
      {
        id: '42->23',
        data: [
          { x: 11, y: 12 },
          { x: 12, y: 12 },
        ],
      },
    ];
    jest.spyOn(React, 'useState').mockImplementationOnce(() => [expectedData, jest.fn()]);
    const wrapper = shallow(<App />);

    expect(wrapper.find('ResponsiveLine').prop('data')).toStrictEqual(expectedData);
  });

  it('should request data from server every 15 seconds', async () => {
    const setStateSpy = jest.fn();
    jest.spyOn(React, 'useState').mockImplementationOnce(() => [[], setStateSpy]);

    mockReturnData = [{ id: '23->42', data: [{ x: 1, y: 2 }] }];
    await act(async () => {
      render(<App />);
    });
    // first request on mount
    await waitFor(async () => expect(setStateSpy).toHaveBeenCalledWith(mockReturnData));

    mockReturnData = [
      { id: '23->42', data: [{ x: 1, y: 2 }] },
      {
        id: '42->23',
        data: [
          { x: 11, y: 12 },
          { x: 12, y: 12 },
        ],
      },
    ];
    await act(async () => {
      // another request after 15 seconds
      jest.advanceTimersByTime(15001);
    });
    await waitFor(async () => expect(setStateSpy).toHaveBeenCalledWith(mockReturnData));

    expect(global.fetch).toBeCalledWith('http://localhost:8080/journeys', {
      method: 'GET',
      headers: { Accept: 'application/json; version=1' },
    });
  });

  it('should cleanup setInterval on unmount', async () => {
    jest.spyOn(global, 'clearInterval');
    await act(async () => {
      render(<App />);
    });

    expect(global.clearInterval).not.toHaveBeenCalled();

    // clearInterval on unmount
    cleanup();
    expect(global.clearInterval).toHaveBeenCalled();
  });
});
