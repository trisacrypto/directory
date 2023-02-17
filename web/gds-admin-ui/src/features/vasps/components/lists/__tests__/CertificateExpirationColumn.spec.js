import dayjs from 'dayjs';
import { faker } from '@faker-js/faker';

import { render, screen } from '@/utils/test-utils';

import CertificateExpirationColumn from '../CertificateExpirationColumn';

describe('<CertificateExpirationColumn />', () => {
  it('should display data correctly', () => {
    const row = {
      original: {
        id: faker.datatype.uuid(),
        certificate_expiration: faker.date.recent(),
      },
    };
    render(<CertificateExpirationColumn row={row} />);

    expect(screen.getByTestId('certificate_expiration').textContent).toBe(
      dayjs(row?.original?.certificate_expiration).format('MMM DD, YYYY h:mm:ss a')
    );
  });

  it('should display N/A', () => {
    const row = {
      original: {
        id: faker.datatype.uuid(),
        certificate_expiration: null,
      },
    };
    render(<CertificateExpirationColumn row={row} />);

    expect(screen.getByTestId('certificate_expiration').textContent).toBe('N/A');
  });
});
