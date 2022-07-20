import React from "react";
import { IconButton, Menu, MenuButton, MenuItem, MenuList } from "@chakra-ui/react";
import { BsThreeDots } from "react-icons/bs";

export type TMenuItem = {
    label: string;
    onClick?: (event: React.MouseEvent<HTMLElement>) => void
};

type OptionsMenuProps = {
    menuItems: TMenuItem[]
};

const OptionsMenu = ({ menuItems }: OptionsMenuProps) => (<Menu>
    <MenuButton
        as={IconButton}
        icon={<BsThreeDots />}
        background="transparent"
        _active={{ outline: 'none' }}
        _focus={{ outline: 'none' }}
        borderRadius={50}
    />
    <MenuList>
        {
            menuItems.map(item => (
                <MenuItem key={item.label} onClick={item.onClick}>{item.label}</MenuItem>
            ))
        }
    </MenuList>
</Menu>);

export default OptionsMenu;
