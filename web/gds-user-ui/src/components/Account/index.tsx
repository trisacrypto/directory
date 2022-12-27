import { Avatar, HStack, Text, Box } from '@chakra-ui/react';
import { Organization } from 'modules/dashboard/organization/organizationType';
import { BsChevronRight } from 'react-icons/bs';
import { Link } from 'react-router-dom';
import { APP_PATH } from 'utils/constants';
type AccountProps = {
  src?: string;
  alt?: string;
  isCurrent: boolean;
} & Partial<Organization>;

export const Account = ({ domain, name, src, id, isCurrent }: AccountProps) => {
  const orgLink = `${APP_PATH.SWITCH_ORGANIZATION}/${id}?vaspName=${name}&vaspDomain=${domain}`;

  const selectOrgHandler = () => {
    window.location.href = orgLink;
  };
  return (
    <>
      <Link onClick={selectOrgHandler} to={orgLink} role="group">
        <HStack
          direction="row"
          w="100%"
          _groupHover={{ bg: '#BEE3F81d' }}
          padding={3}
          justifyContent={'space-between'}
          transition="background 400ms">
          <HStack width="100%" alignItems="center">
            <HStack>
              <Avatar bg="blue" name={name} src={src} fontWeight={700} />
              <HStack>
                <Box>
                  <Text fontWeight={700} data-testid="vaspName">
                    {name}
                  </Text>
                  <Text data-testid="vaspDomain">{domain}</Text>
                </Box>
              </HStack>
            </HStack>
          </HStack>
          {isCurrent && (
            <Text fontStyle={'italic'} fontSize="sm">
              Current{' '}
            </Text>
          )}
          <BsChevronRight />
        </HStack>
      </Link>
    </>
  );
};
