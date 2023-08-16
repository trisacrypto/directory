import {
  Stack,
  Box,
  Table,
  TableCaption,
  Heading,
  Thead,
  Tr,
  Th,
  Tbody,
  Td,
  Badge,
  IconButton
} from '@chakra-ui/react';
import { t, Trans } from '@lingui/macro';
import ConfirmIdentityCertificateModal from 'components/ConfirmIdentityCertificateRequest';
import CkLazyLoadImage from 'components/LazyImage';
import { NoData } from 'components/NoData';
import FormLayout from 'layouts/FormLayout';
import { Certificate } from 'types/type';
import PolygonIcon from 'assets/polygon.svg';
import { useNavigate } from 'react-router-dom';
import formatDisplayedDate from 'utils/formatDisplayedDate';

type MainnetCertificatesProps = {
  data: Certificate[];
  network: 'mainnet' | 'testnet';
};

function MainnetTestnetCertificates({ network, data }: MainnetCertificatesProps) {
  const navigate = useNavigate();
  const handleDetailsClick = (certificateId: string) => {
    navigate(`/dashboard/certificate-inventory/${certificateId}?network=${network}`);
  };

  return (
    <Stack spacing={5}>
      <Box>
        <FormLayout overflowX={'scroll'}>
          <Table variant="simple">
            <TableCaption placement="top" textAlign="start" p={0} m={0} marginBottom={10}>
              <Stack direction={'row'} alignItems={'center'} justifyContent={'space-between'}>
                <Heading fontSize={'1.2rem'} data-testid="title">
                  {network === 'mainnet' ? (
                    <Trans>MainNet Identity Certificates</Trans>
                  ) : (
                    <Trans>TestNet Identity Certificates</Trans>
                  )}
                </Heading>
                <ConfirmIdentityCertificateModal marginLeft="auto !important">
                  <Trans>Request New Identity Certificate</Trans>
                </ConfirmIdentityCertificateModal>
              </Stack>
            </TableCaption>
            <Thead>
              <Tr>
                <Th>
                  <Trans>Serial No</Trans>
                </Th>
                <Th>
                  <Trans>Issue Date</Trans>
                </Th>
                <Th>
                  <Trans>Expiration Date</Trans>
                </Th>
                <Th>
                  <Trans>Status</Trans>
                </Th>
                <Th textAlign="center">
                  <Trans>Actions</Trans>
                </Th>
              </Tr>
            </Thead>
            <Tbody paddingTop={8}>
              {data && data.length ? (
                data?.map((certificate, idx) => (
                  <Tr key={idx} data-testid="table-row">
                    <Td>{certificate.serial_number}</Td>
                    <Td data-testid="issued_at">{formatDisplayedDate(certificate?.issued_at)}</Td>
                    <Td data-testid="expired_at">{formatDisplayedDate(certificate?.expires_at)}</Td>
                    <Td>
                      {certificate.revoked ? (
                        <Badge
                          colorScheme="green"
                          borderRadius="xl"
                          fontWeight={600}
                          textTransform="capitalize"
                          data-testid="revoked">
                          <Trans>Active</Trans>
                        </Badge>
                      ) : (
                        <Badge
                          colorScheme="red"
                          borderRadius="xl"
                          fontWeight={600}
                          textTransform="capitalize"
                          data-testid="revoked">
                          <Trans>Expired</Trans>
                        </Badge>
                      )}
                    </Td>
                    <Td textAlign="center">
                      <IconButton
                        onClick={() => handleDetailsClick(certificate?.serial_number)}
                        aria-label={t`view details`}
                        bg="transparent"
                        _hover={{ bg: 'transparent' }}
                        data-testid="details_btn">
                        <CkLazyLoadImage mx="auto" src={PolygonIcon} />
                      </IconButton>
                    </Td>
                  </Tr>
                ))
              ) : (
                <Tr>
                  <Td colSpan={6} data-testid="no-data">
                    <NoData label={<Trans>No Certificate available</Trans>} />
                  </Td>
                </Tr>
              )}
            </Tbody>
          </Table>
        </FormLayout>
      </Box>
    </Stack>
  );
}

export default MainnetTestnetCertificates;
