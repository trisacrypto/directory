import { render, screen } from '@testing-library/react';

import EmailLog from '../EmailLog';

describe('EmailLogo', () => {
  let emailLogMock;
  it('Should show email log data', () => {
    emailLogMock = [
      {
        contact: 'technical',
        reason: 'verify_contact',
        subject: 'TRISA: Please verify your email address',
        timestamp: '2021-11-08T20:28:36Z',
      },
    ];

    const { container } = render(<EmailLog data={emailLogMock} />);

    expect(container.childNodes.length).toEqual(1);

    expect(screen.getByTestId('email-log-subject').textContent).toBe(emailLogMock[0].subject);
    expect(screen.getByTestId('email-log-subject-line').textContent).toBe(
      'subject: TRISA: Please verify your email address'
    );
    expect(screen.getByTestId('email-log-reason').textContent).toBe(emailLogMock[0].reason);
    expect(screen.getByTestId('email-log-contact').textContent).toBe(emailLogMock[0].contact);
    expect(screen.getByTestId('email-log-contact-line').textContent).toBe('contact: technical');
    expect(screen.getByTestId('email-log-timestamp').textContent).toBe(
      'Mon, 08 Nov 2021 20:28:36 GMT'
    );
  });

  it('should render no data when there are no email log', () => {
    emailLogMock = [];
    render(<EmailLog data={emailLogMock} />);

    const noData = screen.getByTestId('no-data');
    expect(noData).toBeInTheDocument();
    expect(noData.textContent).toBe('No Email Log');
  });
});
