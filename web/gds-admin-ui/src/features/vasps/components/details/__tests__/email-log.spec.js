import { render, screen } from '@testing-library/react';

import EmailLog from '../EmailLog';

describe('EmailLogo', () => {
    let emailLogMock;
    it('Should display email log data', () => {
        emailLogMock = [
            {
                recipient: 'technical',
                reason: 'verify_contact',
                subject: 'TRISA: Please verify your email address',
                timestamp: '2021-11-08T20:28:36Z',
            },
        ];

        const { container } = render(<EmailLog data={emailLogMock} />);

        expect(container.childNodes.length).toEqual(1);

        expect(screen.getByTestId('email-log-subject').textContent).toBe(emailLogMock[0].subject);
        expect(screen.getByTestId('email-log-subject-line').textContent).toBe(
            'Subject: TRISA: Please verify your email address'
        );
        expect(screen.getByTestId('email-log-reason').textContent).toBe(emailLogMock[0].reason);
        expect(screen.getByTestId('email-log-contact').textContent).toBe(emailLogMock[0].recipient);
        expect(screen.getByTestId('email-log-contact-line').textContent).toBe('Recipient: technical');
        expect(screen.getByTestId('email-log-timestamp').textContent).toBe('Date: Mon, 08 Nov 2021 20:28:36 GMT');
    });

    it('should display fallback text', () => {
        emailLogMock = [
            {
                recipient: null,
                reason: null,
                subject: null,
                timestamp: '2021-11-08T20:28:36Z',
            },
        ];

        const { container } = render(<EmailLog data={emailLogMock} />);

        expect(container.childNodes.length).toEqual(1);

        expect(screen.getByTestId('email-log-subject').textContent).toBe('No subject');
        expect(screen.getByTestId('email-log-reason').textContent).toBe('N/A');
        expect(screen.getByTestId('email-log-contact').textContent).toBe('No recipient');
    });

    it('should render no data when there are no email log', () => {
        emailLogMock = [];
        render(<EmailLog data={emailLogMock} />);

        const noData = screen.getByTestId('no-data');
        expect(noData).toBeInTheDocument();
        expect(noData.textContent).toBe('No Email Log');
    });
});
