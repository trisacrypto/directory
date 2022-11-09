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
import Statistics from './Statistics';
import PolygonIcon from 'assets/polygon.svg';
import dayjs from 'dayjs';

type MainnetCertificatesProps = {
  data: Certificate[];
  network: 'mainnet' | 'testnet';
};

const DATE_FORMAT = 'DD-MM-YYYY';

function MainnetTestnetCertificates({ network, data }: MainnetCertificatesProps) {
  const handleDetailsClick = () => {
    // eslint-disable-next-line no-warning-comments
    // TODO: navigate the user to certificate details page
  };

  return (
    <Stack spacing={5}>
      <Statistics />
      <Box>
        <FormLayout overflowX={'scroll'}>
          <Table variant="simple">
            <TableCaption placement="top" textAlign="start" p={0} m={0}>
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
            <Tbody>
              {data && data.length ? (
                data?.map((certificate, idx) => (
                  <Tr key={idx} data-testid="table-row">
                    <Td>{certificate.serial_number}</Td>
                    <Td>{dayjs(certificate.issued_at).format(DATE_FORMAT)}</Td>
                    <Td>{dayjs(certificate.expires_at).format(DATE_FORMAT)}</Td>
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
                        onClick={handleDetailsClick}
                        aria-label={t`view details`}
                        bg="transparent"
                        _hover={{ bg: 'transparent' }}>
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
