import React, { FC, useEffect } from 'react';
import {
  Stack,
  Box,
  Text,
  Heading,
  Table,
  Tbody,
  Tr,
  Td,
  Button,
  Tag,
  Link
} from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { getStepData } from 'utils/utils';
import { loadDefaultValueFromLocalStorage, TStep } from 'utils/localStorageHelper';
import useCertificateStepper from 'hooks/useCertificateStepper';
import {
  BUSINESS_CATEGORY,
  getBusinessCategiryLabel,
  vaspCategories
} from 'constants/basic-details';
interface BasicDetailsReviewProps {}

const BasicDetailsReview = (props: BasicDetailsReviewProps) => {
  const { jumpToStep } = useCertificateStepper();
  const steps: TStep[] = useSelector((state: RootStateOrAny) => state.stepper.steps);
  const [basicDetail, setBasicDetail] = React.useState<any>({});

  useEffect(() => {
    const getStepperData = loadDefaultValueFromLocalStorage();
    const stepData = {
      website: getStepperData.website,
      established_on: getStepperData.established_on,
      vasp_categories: getStepperData.vasp_categories,
      business_category: getStepperData.business_category
    };

    setBasicDetail(stepData);
  }, [steps]);
  return (
    <Box
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      bg={'white'}
      color={'#252733'}
      maxHeight={367}
      fontSize={18}
      p={5}
      px={5}>
      <Stack width={'100%'}>
        <Box display={'flex'} justifyContent="space-between" pt={4} ml={0}>
          <Heading fontSize={20} mb="2rem">
            Section 1: Basic Details
          </Heading>
          <Button
            bg={colors.system.blue}
            color={'white'}
            onClick={() => jumpToStep(1)}
            height={'34px'}
            _hover={{
              bg: '#10aaed'
            }}>
            Edit
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
                <Td borderBottom={'none'}>Website</Td>
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
                <Td>Business Category</Td>
                <Td>{(BUSINESS_CATEGORY as any)[basicDetail.business_category] || 'N/A'}</Td>
                <Td></Td>
              </Tr>
              <Tr borderStyle={'hidden'}>
                <Td whiteSpace="break-spaces" lineHeight={1.5}>
                  Date of Incorporation/ Establishment
                </Td>
                <Td>{basicDetail.established_on || 'N/A'}</Td>
                <Td></Td>
              </Tr>
              <Tr borderStyle={'hidden'}>
                <Td whiteSpace="break-spaces" lineHeight={1.5}>
                  VASP Category
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
