// // Libraries
// import React from "react";
// import { createBrowserRouter, createRoutesFromElements, Route} from "react-router-dom";

// // Components
// import Login from "../Auth/forms/Login";
// import Signup from "../Auth/forms/Signup";

// import { Home } from "../_root/pages";
// import AuthLayout from "../Auth/AuthLayout";
// import RootLayout from "../_root/RootLayout";

// const routes = createBrowserRouter(
//   createRoutesFromElements(
//     <Route>
//       {/* public routes */}
//       <Route element={<AuthLayout />}>
//         <Route path="/login" element={<Login />} />
//         <Route path="/signup" element={<Signup />} />
//       </Route>
      
//       {/* private routes */}
//       <Route element={<RootLayout />}>
//         <Route index element={<Home />} />
//       </Route>
//     </Route>
//   )
// )

// export default routes;