import TableRow from '../TableRow';
import OptionsMenu, { TMenuItem } from '../OptionsMenu';
import React from 'react';
import { t } from '@lingui/macro';

const Menu = () => {
  const MENU_ITEMS: TMenuItem[] = [
    {
      label: t`Edit`
    },
    {
      label: t`Change Permissions`
    },
    {
      label: t`Deactivate`
    }
  ];

  return <OptionsMenu menuItems={MENU_ITEMS} />;
};

const rows = [
  {
    id: '18001',
    signatureId: 'Jones Ferdinand',
    expirationDate: '14/01/2022',
    issueDate: '14/01/2022',
    status: 'active',
    options: <Menu />
  },
  {
    id: '18002',
    signatureId: 'Jones Ferdinand',
    expirationDate: '14/01/2022',
    issueDate: '14/01/2022',
    status: 'active',
    options: <Menu />
  },
  {
    id: '18003',
    signatureId: 'Jones Ferdinand',
    expirationDate: '14/01/2022',
    issueDate: '14/01/2022',
    status: 'active',
    options: <Menu />
  }
];

const SealingCertificateTableRows: React.FC = () => {
  return (
    <>
      {rows.map((row) => (
        <TableRow key={row.id} row={row} />
      ))}
    </>
  );
};

export default SealingCertificateTableRows;
