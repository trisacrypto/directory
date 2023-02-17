import { render, screen } from '@testing-library/react';

import AuditLog from '../AuditLog';

describe('AuditLogo', () => {
  let auditLogMock;

  it('should render', () => {
    auditLogMock = [];
    const { container } = render(<AuditLog data={auditLogMock} />);
    expect(container.childElementCount).toEqual(1);
  });

  it('Should show audit log description', () => {
    auditLogMock = [
      {
        timestamp: '2021-06-17T18:37:23Z',
        previous_state: 'SUBMITTED',
        current_state: 'SUBMITTED',
        description: 'description',
        source: 'source',
      },
    ];

    render(<AuditLog data={auditLogMock} />);
    const auditLogDesc = screen.getByTestId('audit-log-desc');
    expect(auditLogDesc).toBeInTheDocument();
    expect(auditLogDesc.textContent).toBe(auditLogMock[0].description);
  });

  it('Should show formated state', () => {
    auditLogMock = [
      {
        timestamp: '2021-06-17T18:37:23Z',
        previous_state: 'PENDING_REVIEW',
        current_state: 'NO_VERIFICATION',
        description: 'description',
        source: 'source',
      },
    ];

    render(<AuditLog data={auditLogMock} />);
    const auditLogState = screen.getByTestId('audit-log-state');
    const textContent = 'source from Pending Review → No Verification';
    expect(auditLogState).toBeInTheDocument();
    expect(auditLogState.textContent).toBe(textContent);
  });
});
