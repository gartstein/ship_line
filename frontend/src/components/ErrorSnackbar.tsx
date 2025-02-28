// src/components/ErrorSnackbar.tsx
import React from "react";
import { Snackbar, Alert } from "@mui/material";

const ErrorSnackbar = ({ error, setError }: { error: string; setError: (msg: string) => void }) => {
    return (
        <Snackbar open={!!error} autoHideDuration={4000} onClose={() => setError("")}>
            {error ? (
                <Alert severity="error" onClose={() => setError("")}>
                    {error}
                </Alert>
            ) : undefined}
        </Snackbar>
    );
};

export default ErrorSnackbar;