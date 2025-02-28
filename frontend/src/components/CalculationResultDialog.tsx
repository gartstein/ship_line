// src/components/CalculationResultDialog.tsx
import React from "react";
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, Typography } from "@mui/material";

interface CalculationResultDialogProps {
    open: boolean;
    onClose: () => void;
    result: { pack: number; count: number }[];
    itemsOrdered: number | null;
    totalItemsUsed: number | null;
}

const CalculationResultDialog: React.FC<CalculationResultDialogProps> = ({ open, onClose, result, itemsOrdered, totalItemsUsed }) => {
    return (
        <Dialog open={open} onClose={onClose}>
            <DialogTitle>Calculation Result</DialogTitle>
            <DialogContent>
                {itemsOrdered !== null && totalItemsUsed !== null && (
                    <Typography>
                        <strong>Items Ordered:</strong> {itemsOrdered} <br />
                        <strong>Total Items Used:</strong> {totalItemsUsed}
                    </Typography>
                )}
                <Typography variant="h6" sx={{ mt: 2 }}>
                    Pack Breakdown:
                </Typography>
                {result.length > 0 ? (
                    result.map((r, index) => (
                        <Typography key={index}>
                            {r.count} × {r.pack} pack(s)
                        </Typography>
                    ))
                ) : (
                    <Typography>No packs used.</Typography>
                )}
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} color="primary">
                    Close
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default CalculationResultDialog;