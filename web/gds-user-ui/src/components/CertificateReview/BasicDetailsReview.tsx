import React, { useEffect } from 'react';
import {
  Stack,
  Box,
  Heading,
  Table,
  Tbody,
  Tr,
  Td,
  Button,
  Tag,
  Link,
  useColorModeValue
} from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useSelector, RootStateOrAny } from 'react-redux';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { BUSINESS_CATEGORY, getBusinessCategiryLabel } from 'constants/basic-details';
import { Trans } from '@lingui/react';
interface BasicDetailsReviewProps {}

const BasicDetailsReview = (props: BasicDetailsReviewProps) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [basicDetail, setBasicDetail] = React.useState<any>({});
  const textColor = useColorModeValue('gray.800', '#F7F8FC');

  useEffect(() => {
    const getStepperData = loadDefaultValueFromLocalStorage();
    const stepData = {
      website: getStepperData.website,
      established_on: getStepperData.established_on,
      vasp_categories: getStepperData.vasp_categories,
      business_category: getStepperData.business_category
    };
    // '#252733'
    setBasicDetail(stepData);
  }, [steps]);
  return (
    <Box
      border="2px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={textColor}
      maxHeight={367}
      fontSize={18}
      p={5}
      px={5}>
      <Stack width={'100%'}>
        <Box display={'flex'} justifyContent="space-between" pt={4} ml={0}>
          <Heading fontSize={20} mb="2rem">
            <Trans id="Section 1: Basic Details">Section 1: Basic Details</Trans>
          </Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            onClick={() => jumpToStep(1)}
            height={'34px'}
            _hover={{
              bg: '#10aaed'
            }}>
            <Trans id="Edit">Edit</Trans>
          </Button>
        </Box>
        <Stack fontSize={18}>
          <Table
            sx={{
              'td:nth-child(2),td:nth-child(3)': { fontWeight: 'semibold' },
              'td:first-child': {
                width: '50%'
              },
              td: {
                borderBottom: 'none',
                paddingInlineStart: 0,
                paddingY: 2.5
              }
            }}>
            <Tbody
              sx={{
                '*': {
                  fontSize: '1rem'
                }
              }}>
              <Tr>
                <Td borderBottom={'none'} pl={'1rem !important'}>
                  <Trans id="Website">Website</Trans>
                </Td>
                <Td borderBottom={'none'} whiteSpace="break-spaces" lineHeight={1.5}>
                  {basicDetail.website ? (
                    <Link href={basicDetail.website} isExternal>
                      {basicDetail.website}
                    </Link>
                  ) : (
                    'N/A'
                  )}
                </Td>
                <Td></Td>
              </Tr>
              <Tr>
                <Td pl={'1rem !important'}>
                  <Trans id="Business Category">Business Category</Trans>
                </Td>
                <Td>{(BUSINESS_CATEGORY as any)[basicDetail.business_category] || 'N/A'}</Td>
                <Td></Td>
              </Tr>
              <Tr borderStyle={'hidden'}>
                <Td pl={'1rem !important'} whiteSpace="break-spaces" lineHeight={1.5}>
                  <Trans id="Date of Incorporation / Establishment">
                    Date of Incorporation / Establishment
                  </Trans>
                </Td>
                <Td>{basicDetail.established_on || 'N/A'}</Td>
                <Td></Td>
              </Tr>
              <Tr borderStyle={'hidden'}>
                <Td pl={'1rem !important'} whiteSpace="break-spaces" lineHeight={1.5}>
                  <Trans id="VASP Category">VASP Category</Trans>
                </Td>
                <Td>
                  {basicDetail?.vasp_categories && basicDetail?.vasp_categories.length
                    ? basicDetail?.vasp_categories?.map((categ: any) => {
                        return (
                          <Tag key={categ} color={'white'} bg={'blue.400'} mr={2} mb={1}>
                            {getBusinessCategiryLabel(categ)}
                          </Tag>
                        );
                      })
                    : 'N/A'}
                </Td>
                <Td></Td>
              </Tr>
            </Tbody>
          </Table>
        </Stack>
      </Stack>
    </Box>
  );
};
BasicDetailsReview.defaultProps = {
  data: {}
};
export default BasicDetailsReview;
