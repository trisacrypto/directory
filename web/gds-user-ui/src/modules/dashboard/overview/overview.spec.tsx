import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import Overview from '.';
import * as service from './service';

const mockedGetMetrics = jest.spyOn(service, 'getMetrics');

describe('<Overview />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should', () => {
    render(<Overview />);
  });
});
