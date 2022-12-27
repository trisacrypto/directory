import { CopyIcon } from '@chakra-ui/icons';
import {
  Badge,
  Button,
  chakra,
  Heading,
  HStack,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Stack,
  Text,
  VStack
} from '@chakra-ui/react';
import { t, Trans } from '@lingui/macro';
import { getCertificates } from 'application/services/certificates';
import { FileCard } from 'components/FileCard';
import useCertificateStepper from 'hooks/useCertificateStepper';
import FormLayout from 'layouts/FormLayout';
import { isArray } from 'lodash';
import { useEffect, useState } from 'react';
import { HiOutlineDotsHorizontal } from 'react-icons/hi';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { useAsync } from 'react-use';
import { Certificate } from 'types/type';
import downloadFile from 'utils/downloadFile';

type Params = {
  certificateId: string;
};

const formatDisplayedValue = (displayedValue?: string[] | string) => {
  if (isArray(displayedValue)) {
    return displayedValue.join(',');
  }

  return displayedValue ? displayedValue : 'N/A';
};

function CertificateDetails() {
  const [certificateDetails, setCertificateDetails] = useState<Certificate | null>(null);
  const navigate = useNavigate();
  const { jumpToLastStep } = useCertificateStepper();
  const params = useParams<Params>();
  const [searchParams] = useSearchParams();
  const { value } = useAsync(getCertificates);

  const certificateId = params?.certificateId;
  const network = searchParams.get('network') as any;

  useEffect(() => {
    if (certificateId && network && value) {
      const details = value[network]?.find(
        (certificate: Certificate) => certificate?.serial_number === certificateId
      );

      setCertificateDetails(details);
    }
  }, [certificateId, network, value]);

  const handleEditClick = () => {
    navigate('/dashboard/certificate/registration');
    jumpToLastStep();
  };

  const handlePublicIdentityKeyDownloadClick = (data: string) => {
    const filename = 'public-identity-key.pem';
    const mimetype = 'application/x-pem-file';
    downloadFile(data, filename, mimetype);
  };

  const handleTrustChainDownloadClick = (chain: string) => {
    const filename = 'trust-chain-certificate.gz';
    const mimetype = 'application/x-x509-ca-cert';
    downloadFile(chain, filename, mimetype);
  };

  return (
    <>
      <Heading marginBottom="30px">
        <Trans>X.509 Identity Certificate Details</Trans>
      </Heading>
      <Stack spacing={5}>
        <HStack
          border="1px solid #DFE0EB"
          bg="#D8EAF6"
          py={3}
          px={8}
          borderRadius="10px"
          justifyContent="space-between">
          <Text fontWeight={700} fontSize="lg">
            <Trans>MAINNET Identity Certificate</Trans>
          </Text>
          <Button variant="primary" w="100%" maxW="130px" onClick={handleEditClick}>
            <Trans>Edit</Trans>
          </Button>
        </HStack>
        <FormLayout>
          <HStack justifyContent="space-between" w="100%">
            <Text textTransform="capitalize" fontWeight={700} fontSize="lg">
              <Trans>MainNet Certificate Details</Trans>
            </Text>
            <Menu>
              <MenuButton
                variant="ghost"
                color="#858585"
                as={Button}
                margin={0}
                marginInline={0}
                sx={{
                  '& .chakra-button__icon': {
                    marginInline: 0
                  }
                }}
                rightIcon={<HiOutlineDotsHorizontal size={25} />}
              />
              <MenuList>
                <MenuItem>
                  <CopyIcon mr={2} color="blue" /> Copy signature
                </MenuItem>
                <MenuItem>
                  <CopyIcon mr={2} color="blue" /> Copy serial number
                </MenuItem>
              </MenuList>
            </Menu>
          </HStack>
          <VStack align="start">
            <Text>
              <chakra.span fontWeight={700}>
                <Trans>Status</Trans>:
              </chakra.span>{' '}
              {certificateDetails?.revoked ? (
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
            </Text>
            <Text>
              <chakra.span fontWeight={700}>
                <Trans>Serial Number</Trans>:
              </chakra.span>{' '}
              {certificateDetails?.serial_number || 'N/A'}
            </Text>
            <Text>
              <>
                <chakra.span fontWeight={700}>
                  <Trans>Expires</Trans>:
                </chakra.span>{' '}
                {certificateDetails?.expires_at || 'N/A'}
              </>
            </Text>
            <Text>
              <chakra.span fontWeight={700}>
                <Trans>Issuer</Trans>:
              </chakra.span>{' '}
              {certificateDetails?.details.issuer.common_name || 'N/A'}
            </Text>
            <Text>
              <chakra.span fontWeight={700}>
                <Trans>Subject</Trans>:
              </chakra.span>{' '}
              {certificateDetails?.details?.subject?.common_name || 'N/A'}
            </Text>
            <Text>
              <chakra.span fontWeight={700}>
                <Trans>Endpoint</Trans>:
              </chakra.span>{' '}
              {certificateDetails?.details?.endpoint || 'N/A'}
            </Text>
          </VStack>
        </FormLayout>
        <FormLayout>
          <Text textTransform="capitalize" fontWeight={700}>
            <Trans>Download</Trans>
          </Text>

          <Stack direction="row" justifyContent="start" w="100%" spacing={10}>
            <FileCard
              name={t`Public Identity Key`}
              file={certificateDetails?.details?.data}
              ext={`.PEM`}
              onDownload={() =>
                certificateDetails?.details?.data &&
                handlePublicIdentityKeyDownloadClick(certificateDetails?.details?.data)
              }
            />
            <FileCard
              file={certificateDetails?.details?.chain}
              name="TRISA Trust Chain (CA)"
              ext=".GZ"
              onDownload={() =>
                certificateDetails?.details?.chain &&
                handleTrustChainDownloadClick(certificateDetails?.details?.chain)
              }
            />
          </Stack>
        </FormLayout>
        <Stack direction="row" justifyContent="space-between" spacing={10}>
          <FormLayout w="100%">
            <HStack justifyContent="space-between" w="100%">
              <Text textTransform="capitalize" fontWeight={700} fontSize="lg">
                <Trans>User Details</Trans>
              </Text>
            </HStack>
            <VStack align="start">
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Common Name</Trans>:
                </chakra.span>{' '}
                {certificateDetails?.details?.issuer?.common_name || 'N/A'}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Country</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.issuer?.country)}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Locality</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.issuer?.locality)}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Organization</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.issuer?.organization)}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Organization Unit</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.issuer?.organizational_unit)}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Postal Code</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.issuer?.postal_code)}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Province</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.issuer?.province)}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Serial Number</Trans>:
                </chakra.span>{' '}
                {certificateDetails?.details?.issuer?.serial_number}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Street Address</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.issuer?.street_address)}
              </Text>
            </VStack>
          </FormLayout>
          <FormLayout w="100%">
            <HStack justifyContent="space-between" w="100%">
              <Text textTransform="capitalize" fontWeight={700} fontSize="lg">
                <Trans>Subject Details</Trans>
              </Text>
            </HStack>
            <VStack align="start">
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Common Name</Trans>:
                </chakra.span>{' '}
                {certificateDetails?.details?.subject?.common_name || 'N/A'}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Country</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.subject?.country)}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Locality</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.subject?.locality)}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Organization</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.subject?.organization)}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Organization Unit</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.subject?.organizational_unit)}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Postal Code</Trans>:
                </chakra.span>{' '}
                N/A
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Province</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.subject?.province)}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Serial Number</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.subject?.serial_number)}
              </Text>
              <Text>
                <chakra.span fontWeight={700}>
                  <Trans>Street Address</Trans>:
                </chakra.span>{' '}
                {formatDisplayedValue(certificateDetails?.details?.subject?.street_address)}
              </Text>
            </VStack>
          </FormLayout>
        </Stack>
      </Stack>
    </>
  );
}

export default CertificateDetails;
