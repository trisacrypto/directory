import {
  Flex,
  FlexProps,
  HStack,
  IconButton,
  useColorModeValue,
  Text,
  Divider,
  Menu,
  MenuButton,
  Box,
  Avatar,
  MenuList,
  MenuItem,
  MenuDivider,
} from "@chakra-ui/react";
import { FiBell, FiMenu, FiSearch } from "react-icons/fi";

interface MobileProps extends FlexProps {
  onOpen: () => void;
}
const MobileNav = ({ onOpen, ...rest }: MobileProps) => {
  return (
    <Flex
      ml={{ base: 0, md: 60 }}
      px={{ base: 4, md: 4 }}
      height="20"
      alignItems="center"
      bg={useColorModeValue("white", "gray.900")}
      borderBottomWidth="1px"
      borderBottomColor={useColorModeValue("gray.200", "gray.700")}
      justifyContent={{ base: "space-between", md: "flex-end" }}
      {...rest}
    >
      <IconButton
        display={{ base: "flex", md: "none" }}
        borderRadius={50}
        onClick={onOpen}
        variant="outline"
        aria-label="open menu"
        icon={<FiMenu />}
      />

      <Text
        display={{ base: "flex", md: "none" }}
        fontSize="2xl"
        fontFamily="monospace"
        fontWeight="bold"
      >
        Logo
      </Text>

      <HStack spacing={{ base: "0", md: "6" }}>
        <HStack>
          <IconButton
            size="lg"
            variant="ghost"
            aria-label="open menu"
            borderRadius={50}
            color="gray.700"
            _focus={{ boxShadow: "none" }}
            icon={<FiSearch />}
          />
          <IconButton
            size="lg"
            variant="ghost"
            aria-label="open menu"
            borderRadius={50}
            color="gray.700"
            _focus={{ boxShadow: "none" }}
            icon={<FiBell />}
          />
        </HStack>
        <Divider orientation="vertical" height={8} />
        <Menu>
          <MenuButton transition="all 0.3s" _focus={{ boxShadow: "none" }}>
            <HStack>
              <Text fontSize="sm" color="blackAlpha.700">
                Jones Ferdinand
              </Text>
              <Box borderRadius="50%" borderWidth={2} padding={0.5}>
                <Avatar
                  size={"md"}
                  height="43.3"
                  w="43.3"
                  src={
                    "https://images.unsplash.com/photo-1619946794135-5bc917a27793?ixlib=rb-0.3.5&q=80&fm=jpg&crop=faces&fit=crop&h=200&w=200&s=b616b2c5b373a80ffc9636ba24f7a4a9"
                  }
                />
              </Box>
            </HStack>
          </MenuButton>
          <MenuList
            bg={useColorModeValue("white", "gray.900")}
            borderColor={useColorModeValue("gray.200", "gray.700")}
          >
            <MenuItem>Profile</MenuItem>
            <MenuDivider />
            <MenuItem>Sign out</MenuItem>
          </MenuList>
        </Menu>
      </HStack>
    </Flex>
  );
};

export default MobileNav;
