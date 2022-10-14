import TableRow from '../TableRow';
import React from 'react';
import PolygonIcon from 'assets/polygon.svg';
import CkLazyLoadImage from 'components/LazyImage';
import { IconButton } from '@chakra-ui/react';
import { t } from '@lingui/macro';

const rows = [
  {
    id: '18001',
    signatureId: 'Jones Ferdinand',
    expirationDate: '14/01/2022',
    issueDate: '14/01/2022',
    status: 'active',
    options: (
      <IconButton aria-label={t`view details`} bg="transparent" _hover={{ bg: 'transparent' }}>
        <CkLazyLoadImage mx="auto" src={PolygonIcon} />
      </IconButton>
    )
  },
  {
    id: '18002',
    signatureId: 'Jones Ferdinand',
    expirationDate: '14/01/2022',
    issueDate: '14/01/2022',
    status: 'active',
    options: (
      <IconButton aria-label={t`view details`} bg="transparent" _hover={{ bg: 'transparent' }}>
        <CkLazyLoadImage mx="auto" src={PolygonIcon} />
      </IconButton>
    )
  },
  {
    id: '18003',
    signatureId: 'Jones Ferdinand',
    expirationDate: '14/01/2022',
    issueDate: '14/01/2022',
    status: 'active',
    options: (
      <IconButton aria-label={t`view details`} bg="transparent" _hover={{ bg: 'transparent' }}>
        <CkLazyLoadImage mx="auto" src={PolygonIcon} />
      </IconButton>
    )
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
