let map = fn(arr, f) {
    let iter = fn(arr, accumulated) {
        if (len(arr) == 0) { 
            accumulated 
        } else {
            iter(rest(arr), push(accumulated, f(first(arr)))); 
        }; 
    }; 
    
    iter(arr, []); 
};

let reduce = fn(arr, initial, f) {
    let iter = fn(arr, result) {
        if (len(arr) == 0) {
            result
        } else {
            iter(rest(arr), f(result, first(arr)));
        }
    };

    iter(arr, initial)
};

let filter = fn(arr, pred) {
    let iter = fn(arr, accumulated) {
        if (len(arr) == 0) {
            accumulated
        } else {
            let elem = first(arr);
            if (pred(elem)) {
                iter(rest(arr), push(accumulated, elem));
            } else {
                iter(rest(arr), accumulated);
            }
        }
    };

    iter(arr, []);
};

let sum = fn(arr) {
    reduce(arr, 0, fn(initial, el) { initial + el });
};
