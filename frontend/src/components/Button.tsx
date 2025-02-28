import React from "react";
import { Button as MuiButton } from "@mui/material";

interface ButtonProps {
    onClick: () => void;
    loading: boolean;
}

const Button: React.FC<ButtonProps> = ({ onClick, loading }) => {
    return (
        <MuiButton
            variant="contained"
            color="primary"
            onClick={onClick}
            disabled={loading}
        >
            {loading ? "Calculating..." : "Calculate"}
        </MuiButton>
    );
};

export default Button;