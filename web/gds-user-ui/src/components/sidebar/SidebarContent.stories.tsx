import { Meta, Story } from "@storybook/react";
import SidebarContent from "./SidebarContent";

interface SidebarProps {
  onClose: () => void;
}

export default {
  title: "components/SideBarContent",
  component: SidebarContent,
} as Meta;

const Template: Story<SidebarProps> = (args) => <SidebarContent {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
