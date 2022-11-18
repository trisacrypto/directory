import { faker } from '@faker-js/faker';
import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import MainnetTestnetCertificates from '../MainnetTestnetCertificates';

describe('<MainnetTestnetCertificates />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });

  it('should display MainNet Identity Certificates as title when network is testnet', () => {
    render(<MainnetTestnetCertificates data={[]} network={'mainnet'} />);

    expect(screen.getByTestId('title').textContent).toBe('MainNet Identity Certificates');
  });

  it('should display TestNet Identity Certificates as title when the network is testnet', () => {
    render(<MainnetTestnetCertificates data={[]} network={'testnet'} />);

    expect(screen.getByTestId('title').textContent).toBe('TestNet Identity Certificates');
  });

  it('should display NoData component when data is empty', () => {
    render(<MainnetTestnetCertificates data={[]} network={'mainnet'} />);

    expect(screen.getByTestId('no-data')).toBeInTheDocument();
  });

  it('should display table rows when data is not empty', () => {
    const mockData = [
      {
        serial_number: '2312834348738913753151',
        issued_at: '2022-11-16T02:05:56.174Z',
        expires_at: '2022-12-16T11:22:46.015Z',
        revoked: true,
        details: 'rzer'
      },
      {
        serial_number: '954687348674867246',
        issued_at: faker.date.soon(),
        expires_at: faker.date.soon(),
        revoked: false,
        details: 'rzer'
      }
    ];
    render(<MainnetTestnetCertificates data={mockData} network={'mainnet'} />);

    expect(screen.getAllByTestId('table-row').length).toBe(2);

    const revoked = screen.getAllByTestId('revoked').map((node) => node.textContent);
    expect(revoked).toStrictEqual(['Active', 'Expired']);
    expect(screen.getAllByTestId('issued_at').map((node) => node.textContent)[0]).toStrictEqual(
      '16-11-2022'
    );
    expect(screen.getAllByTestId('expired_at').map((node) => node.textContent)[0]).toStrictEqual(
      '16-12-2022'
    );
  });

  it('should navigate', () => {
    const mockData = [
      {
        serial_number: '2312834348738913753151',
        issued_at: '2022-11-16T02:05:56.174Z',
        expires_at: '2022-12-16T11:22:46.015Z',
        revoked: true,
        details: 'rzer'
      },
      {
        serial_number: '954687348674867246',
        issued_at: faker.date.soon(),
        expires_at: faker.date.soon(),
        revoked: false,
        details: 'rzer'
      }
    ];

    render(<MainnetTestnetCertificates data={mockData} network={'mainnet'} />);

    expect(location.pathname).toBe('/');

    userEvent.click(screen.getAllByTestId('details_btn')[0]);

    expect(location.pathname).toBe(`/dashboard/certificate-inventory/${mockData[0].serial_number}`);
  });
});
